package grpc

import (
	"testing"
)

func TestGetTicker(t *testing.T) {
	stock := StockInfo{
		Ticker:      "GOOGL",
		Company:     "Alphabet Inc.",
		Description: "Google's parent company, focusing on search, advertising, and tech innovation.",
	}
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
	result := stock.GetCompany()
	if result != "Alphabet Inc." {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Alphabet Inc.")
	}
}
