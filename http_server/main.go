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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := c.SelectTicker(ctx, &pb.TickerRequest{Ticker: ticker})
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

//Frontend Maangement
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := muxtrace.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/stocks", returnAllStocks)

	myRouter.HandleFunc("/stock", createNewStock).Methods("POST")
	myRouter.HandleFunc("/stock/{ticker}", returnSingleStock)
	myRouter.HandleFunc("/delete/{ticker}", deleteStock)

	log.Fatal(http.ListenAndServe(":1313", myRouter))
}

var (
	addr = flag.String("addr", "backend:50051", "the address to connect to")
	c    = pb.NewTickerManagerClient(nil)
)

func main() {
	flag.Parse()
	//set up tracer
	tracer.Start(tracer.WithAgentAddr("datadog-agent:8126"))
	defer tracer.Stop()

	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))

	// set up grpc
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewTickerManagerClient(conn)

	//take http request
	handleRequests()

}
