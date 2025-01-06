package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/sinta-p/tracker-app/http_server/grpc"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// data structure
type Stock struct {
	Ticker      string `json:"ticker"`
	Company     string `json:"company"`
	Description string `json:"desc"`
}

var Stocks []Stock

// GET single stock
func returnSingleStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ticker := strings.ToUpper(vars["ticker"])

	fmt.Println("Endpoint Hit: returnSingleStock - " + ticker)

	// Contact the backend and print out its response.
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	// method 1: original
	reply, err := c.SelectTicker(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not query: %v", err)
	} else {
		log.Printf("query successful: %s", ticker)
	}
	// method 2: high cpu time
	reply, err = c2.SelectTicker2(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not query: %v", err)
	} else {
		log.Printf("query successful: %s", ticker)
	}
	// method 3: high allocation
	reply, err = c3.SelectTicker3(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not query: %v", err)
	} else {
		log.Printf("query successful: %s", ticker)
	}
	// method 4: high heap
	reply, err = c4.SelectTicker4(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not query: %v", err)
	} else {
		log.Printf("query successful: %s", ticker)
	}
	// method 4: mutex operation
	reply, err = c5.SelectTicker5(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not query: %v", err)
	} else {
		log.Printf("query successful: %s", ticker)
	}

	var targetStock Stock
	targetStock.Ticker = reply.GetTicker()
	targetStock.Company = reply.GetCompany()
	targetStock.Description = reply.GetDescription()

	json.NewEncoder(w).Encode(targetStock)
}

// GET all stocks
func returnAllStocks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllStocks")
	//json.NewEncoder(w).Encode(Stocks)
}

// Creating New Stock
func createNewStock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNew")

	reqBody, _ := ioutil.ReadAll(r.Body)

	var stock Stock
	json.Unmarshal(reqBody, &stock)
	fmt.Println()

	// Send request to backend
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := c.InsertTicker(ctx, &pb.StockInfo{Ticker: stock.Ticker, Company: stock.Company, Description: stock.Description})
	if err != nil {
		log.Fatalf("could not insert %s : %v", stock.Ticker, err)
	} else {
		log.Printf("insert successful: %s", stock.Ticker)
	}

	success_status := reply.GetSuccess()

	json.NewEncoder(w).Encode(success_status)
}

// Delete stock
func deleteStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ticker := vars["ticker"]

	// Send request to backend
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := c.DeleteTicker(ctx, &pb.TickerRequest{Ticker: ticker})
	if err != nil {
		log.Fatalf("could not delete %s : %v", ticker, err)
	} else {
		log.Printf("delete successful: %s", ticker)
	}

	success_status := reply.GetSuccess()

	json.NewEncoder(w).Encode(success_status)

}

// Frontend Maangement
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := muxtrace.NewRouter(muxtrace.WithServiceName("tracer-http")).StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/stocks", returnAllStocks)

	myRouter.HandleFunc("/stock", createNewStock).Methods("POST")
	myRouter.HandleFunc("/stock/{ticker}", returnSingleStock)
	myRouter.HandleFunc("/delete/{ticker}", deleteStock)

	log.Fatal(http.ListenAndServe(":1313", myRouter))
}

var (
	// addr  = flag.String("addr", "backend:50051", "the address to connect to")
	// addr2 = flag.String("addr", "backend:50052", "the address to connect to")
	c  = pb.NewTickerManagerClient(nil)
	c2 = pb.NewTickerManager2Client(nil)
	c3 = pb.NewTickerManager3Client(nil)
	c4 = pb.NewTickerManager4Client(nil)
	c5 = pb.NewTickerManager5Client(nil)
)

func main() {
	flag.Parse()
	//set up tracer
	tracer.Start(
		tracer.WithEnv("dev"),
		tracer.WithService("http-client"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithAgentAddr("datadog-agent.datadog-ns.svc.cluster.local:8126"),
	)
	defer tracer.Stop()

	// set up profiler
	err := profiler.Start(
		profiler.WithService("http-client"),
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
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("grpc"))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("grpc"))

	// set up grpc
	conn, err := grpc.Dial("backend-main:50051", grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewTickerManagerClient(conn)

	// set up grpc 2
	conn2, err := grpc.Dial("backend-high-cpu:50052", grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()
	c2 = pb.NewTickerManager2Client(conn2)

	// set up grpc 3
	conn3, err := grpc.Dial("backend-high-alloc:50053", grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn3.Close()
	c3 = pb.NewTickerManager3Client(conn3)

	// set up grpc 4
	conn4, err := grpc.Dial("backend-high-heap:50054", grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn4.Close()
	c4 = pb.NewTickerManager4Client(conn4)

	// set up grpc 5
	conn5, err := grpc.Dial("backend-mutex:50055", grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn5.Close()
	c5 = pb.NewTickerManager5Client(conn5)

	//take http request
	handleRequests()

}
