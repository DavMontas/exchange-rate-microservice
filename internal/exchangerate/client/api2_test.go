package client

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

func TestAPI2Client_FetchRate_Success_EURtoX(t *testing.T) {
	xmlResp := `
<Envelope>
  <Cube>
    <Cube time="2025-07-27">
      <Cube currency="JPY" rate="168.2"/>
      <Cube currency="USD" rate="1.1"/>
    </Cube>
  </Cube>
</Envelope>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xmlResp))
	}))
	defer srv.Close()

	cli := NewAPI2Client(srv.URL)
	pair := domain.CurrencyPair{From: "EUR", To: "JPY", Amount: 2}
	q := cli.FetchRate(context.Background(), pair)

	if q.Err != nil {
		t.Fatalf("expected no error, got %v", q.Err)
	}
	want := 168.2 * 2
	if math.Abs(q.Rate-want) > 1e-6 {
		t.Errorf("expected rate %.2f, got %.2f", want, q.Rate)
	}
}

func TestAPI2Client_FetchRate_Success_XtoEUR(t *testing.T) {
	xmlResp := `
<Envelope>
  <Cube>
    <Cube time="2025-07-27">
      <Cube currency="USD" rate="1.1"/>
      <Cube currency="JPY" rate="168.2"/>
    </Cube>
  </Cube>
</Envelope>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xmlResp))
	}))
	defer srv.Close()

	cli := NewAPI2Client(srv.URL)
	pair := domain.CurrencyPair{From: "USD", To: "EUR", Amount: 5}
	q := cli.FetchRate(context.Background(), pair)

	if q.Err != nil {
		t.Fatalf("expected no error, got %v", q.Err)
	}
	want := (1.0 / 1.1) * 5
	if math.Abs(q.Rate-want) > 1e-6 {
		t.Errorf("expected rate %.6f, got %.6f", want, q.Rate)
	}
}
