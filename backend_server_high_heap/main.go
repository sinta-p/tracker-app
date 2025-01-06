package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sinta-p/tracker-app/backend_server_high_heap/DatabaseManager"
	pb "github.com/sinta-p/tracker-app/backend_server_high_heap/grpc"
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

// func to create a memory leak problem through go runtime

var size = 64 * 1024 //65536

func smallButLongFunction() {
	dataLarge := make([]byte, size)

	//to simulate an expensive task
	time.Sleep(time.Second * 60 * 60) //1 hour

	fmt.Printf("len=%d cap=%d %v\n",
		len(dataLarge), cap(dataLarge), dataLarge)
}

func (s *server) SelectTicker4(ctx context.Context, in *pb.TickerRequest) (*pb.StockInfo, error) {

	stock, err := DatabaseManager.DBSelectTicker(ctx, db, in.GetTicker())
	if err != nil {
		log.Printf("Unexpected error, err:%s", err)
	}

	go smallButLongFunction()

	return &pb.StockInfo{Ticker: stock.Ticker, Company: stock.Company, Description: stock.Description}, err
}

// grpc declarations
var (
	port = flag.Int("port", 50054, "The server port")
)

type server struct {
	pb.UnimplementedTickerManager4Server
}

// Main
func main() {
	flag.Parse()

	//set up tracer
	tracer.Start(
		tracer.WithEnv("dev"),
		tracer.WithService("ticker-manager4"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithAgentAddr("datadog-agent.datadog-ns.svc.cluster.local:8126"),
	)
	defer tracer.Stop()

	// set up profiler
	err := profiler.Start(
		profiler.WithService("ticker-manager4"),
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
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("ticker-manager4"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("ticker-manager4"))

	// Create a listener for the server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterTickerManager4Server(s, &server{})
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
