package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

func TestAPI3Client_FetchRate_Success(t *testing.T) {
	resp := map[string]map[string]float64{"rates": {"EUR": 0.75}}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(resp)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	cli := NewAPI3Client(srv.URL + "?from=%s&to=%s")
	pair := domain.CurrencyPair{From: "USD", To: "EUR", Amount: 300}
	q := cli.FetchRate(context.Background(), pair)

	if q.Err != nil {
		t.Fatalf("expected no error, got %v", q.Err)
	}
	if q.Rate != 0.75*300 {
		t.Errorf("expected rate %.2f, got %.2f", 0.75*300, q.Rate)
	}
}

func TestAPI3Client_FetchRate_NotFound(t *testing.T) {
	empty := map[string]map[string]float64{"rates": {}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(empty)
	}))
	defer srv.Close()

	cli := NewAPI3Client(srv.URL + "?from=%s&to=%s")
	pair := domain.CurrencyPair{From: "USD", To: "XXX", Amount: 10}
	q := cli.FetchRate(context.Background(), pair)
	if q.Err == nil {
		t.Fatal("expected error for missing currency, got none")
	}
}
