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
	pb "github.com/sinta-p/tracker-app/grpc"

	"google.golang.org/grpc"
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

	// Contact the server and print out its response.
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
	json.NewEncoder(w).Encode(Stocks)
}

// Creating New Stock
func createNewStock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNew")

	reqBody, _ := ioutil.ReadAll(r.Body)

	var stock Stock
	json.Unmarshal(reqBody, &stock)
	fmt.Println()
	Stocks = append(Stocks, stock)

	json.NewEncoder(w).Encode(stock)
}

// Delete stock
func deleteStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ticker := vars["ticker"]

	for index, stock := range Stocks {
		if stock.Ticker == ticker {
			Stocks = append(Stocks[:index], Stocks[index+1:]...)
		}
	}
}

//Frontend Maangement
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/stocks", returnAllStocks)

	myRouter.HandleFunc("/stock", createNewStock).Methods("POST")
	myRouter.HandleFunc("/stock/{ticker}", returnSingleStock)

	log.Fatal(http.ListenAndServe(":1313", myRouter))
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	c    = pb.NewTickerManagerClient(nil)
)

func main() {
	flag.Parse()
	// set up grpc
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewTickerManagerClient(conn)

	//take http request
	handleRequests()

}
