package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

func TestAPI1Client_FetchRate_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]struct {
			Rate float64 `json:"rate"`
		}{
			"eur": {Rate: 0.5},
		}
		json.NewEncoder(w).Encode(data)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	cli := NewAPI1Client(srv.URL + "/daily/%s.json")
	pair := domain.CurrencyPair{From: "USD", To: "EUR", Amount: 100}
	q := cli.FetchRate(context.Background(), pair)

	if q.Err != nil {
		t.Fatalf("expected no error, got %v", q.Err)
	}
	want := 0.5 * 100
	if q.Rate != want {
		t.Errorf("expected rate %.2f, got %.2f", want, q.Rate)
	}
}

func TestAPI1Client_FetchRate_NotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]struct {
			Rate float64 `json:"rate"`
		}{}
		json.NewEncoder(w).Encode(data)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	cli := NewAPI1Client(srv.URL + "/daily/%s.json")
	pair := domain.CurrencyPair{From: "USD", To: "XXX", Amount: 1}
	q := cli.FetchRate(context.Background(), pair)

	if q.Err == nil {
		t.Fatal("expected error for missing currency, got none")
	}
}
