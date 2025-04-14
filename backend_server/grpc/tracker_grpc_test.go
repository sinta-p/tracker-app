package grpc

import (
	"testing"
	"time"
)

func TestGetTicker(t *testing.T) {
	stock := StockInfo{
		Ticker:      "GOOGL",
		Company:     "Alphabet Inc.",
		Description: "Google's parent company, focusing on search, advertising, and tech innovation.",
	}
	time.Sleep(10 * time.Second)

	result := stock.GetTicker()
	if result != "GOOGL" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "GOOGL")
	}
}

func TestGetCompany(t *testing.T) {
	stock := StockInfo{
		Ticker:      "GOOGL",
		Company:     "Alphabet Inc.",
		Description: "Google's parent company, focusing on search, advertising, and tech innovation.",
	}
	time.Sleep(10 * time.Second)

	result := stock.GetCompany()
	if result != "Alphabet Inc." {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Alphabet Inc.")
	}
}
