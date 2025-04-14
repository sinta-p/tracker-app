package main

import (
	"testing"

	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

type StockInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ticker      string `protobuf:"bytes,1,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Company     string `protobuf:"bytes,2,opt,name=company,proto3" json:"company,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

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
