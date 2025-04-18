package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/sinta-p/tracker-app/backend_server_mutex/DatabaseManager"
	pb "github.com/sinta-p/tracker-app/backend_server_mutex/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var db = DatabaseManager.OpenDatabase()

func (s *server) InsertTicker(ctx context.Context, in *pb.StockInfo) (*pb.Status, error) {
	table := "stocks_tab"
	columns := "ticker,company,description"
	ticker := in.GetTicker()
	company := in.GetCompany()
	description := in.GetDescription()
	values := "\"" + ticker + "\", \"" + company + "\", \"" + description + "\""

	err := DatabaseManager.DBInsert(ctx, db, table, columns, values)

	if err != nil {
		log.Printf("Unexpected error, err:%s", err)
		return &pb.Status{Success: false}, err
	} else {
		return &pb.Status{Success: true}, err
	}
}

func (s *server) DeleteTicker(ctx context.Context, in *pb.TickerRequest) (*pb.Status, error) {
	table := "stocks_tab"
	column := "ticker"
	value := in.GetTicker()

	err := DatabaseManager.DBDelete(ctx, db, table, column, value)

	if err != nil {
		log.Printf("Unexpected error, err:%s", err)
		return &pb.Status{Success: false}, err
	} else {
		return &pb.Status{Success: true}, err
	}
}

// Counter operations
type Counter struct {
	mu       sync.Mutex
	counter  int
	cooldown int
}

var c = Counter{
	counter:  0,
	cooldown: 1,
}

func (c *Counter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	//error: used counter instead of cooldown leading to constant increase in blocking time
	time.Sleep(time.Duration(c.counter) * time.Second)
	c.counter += 1
	log.Printf("counter value:%d", c.counter)
}

func (s *server) SelectTicker5(ctx context.Context, in *pb.TickerRequest) (*pb.StockInfo, error) {

	stock, err := DatabaseManager.DBSelectTicker(ctx, db, in.GetTicker())
	if err != nil {
		log.Printf("Unexpected error, err:%s", err)
	}

	go c.inc()

	return &pb.StockInfo{Ticker: stock.Ticker, Company: stock.Company, Description: stock.Description}, err
}

// grpc declarations
var (
	port = flag.Int("port", 50055, "The server port")
)

type server struct {
	pb.UnimplementedTickerManager5Server
}

// Main
func main() {
	flag.Parse()

	//set up tracer
	tracer.Start(
		tracer.WithEnv("dev"),
		tracer.WithService("ticker-manager5"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithAgentAddr("datadog-agent.datadog-ns.svc.cluster.local:8126"),
	)
	defer tracer.Stop()

	// set up profiler
	err := profiler.Start(
		profiler.WithService("ticker-manager5"),
		profiler.WithEnv("dev"),
		profiler.WithVersion("1.0.0"),
		profiler.WithTags("owner:sin,app:tracker-app"),
		profiler.WithAgentAddr("datadog-agent.datadog-ns.svc.cluster.local:8126"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			// The profiles below are disabled by default to keep overhead
			// low, but can be enabled as needed.

			profiler.BlockProfile,
			profiler.MutexProfile,
			profiler.GoroutineProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("ticker-manager5"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("ticker-manager5"))

	// Create a listener for the server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterTickerManager5Server(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Testing main
// func main() {
// 	SelectTicker("DDOG")
// 	InsertNewTicker("SE", "SEA Limited", "E-commerce, Fintech & Gaming")
// }
