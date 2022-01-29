package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/sinta-p/tracker-app/backend_server/DatabaseManager"
	pb "github.com/sinta-p/tracker-app/backend_server/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

func (s *server) SelectTicker(ctx context.Context, in *pb.TickerRequest) (*pb.StockInfo, error) {

	stock, err := DatabaseManager.DBSelectTicker(ctx, db, in.GetTicker())
	if err != nil {
		log.Printf("Unexpected error, err:%s", err)
	}

	return &pb.StockInfo{Ticker: stock.Ticker, Company: stock.Company, Description: stock.Description}, err
}

// grpc declarations
var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedTickerManagerServer
}

// Main
func main() {
	flag.Parse()

	//set up tracer
	tracer.Start(tracer.WithAgentAddr("datadog-agent:8126"))
	defer tracer.Stop()

	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("TicketManager"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("TicketManager"))

	// Create a listener for the server
	lis, err := net.Listen("tcp", fmt.Sprintf("backend:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	pb.RegisterTickerManagerServer(s, &server{})
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
