package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/client"
	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type mockClient struct {
	quote domain.Quote
}

func (m *mockClient) FetchRate(_ context.Context, _ domain.CurrencyPair) domain.Quote {
	return m.quote
}

func TestBestQuote_AllSuccess(t *testing.T) {
	clients := []client.RateClient{
		&mockClient{quote: domain.Quote{Provider: "a", Rate: 1.0}},
		&mockClient{quote: domain.Quote{Provider: "b", Rate: 2.0}},
	}
	nopLogger := zap.NewNop().Sugar()
	svc := NewExchangeService(clients, 1*time.Second, nopLogger)
	q := svc.BestQuote(context.Background(), domain.CurrencyPair{})
	if q.Provider != "b" || q.Rate != 2.0 {
		t.Errorf("expected best b with rate 2.0, got %v %.2f", q.Provider, q.Rate)
	}
}

func TestBestQuote_PartialFailures(t *testing.T) {
	clients := []client.RateClient{
		&mockClient{quote: domain.Quote{Provider: "a", Err: errors.New("fail")}},
		&mockClient{quote: domain.Quote{Provider: "b", Rate: 3.0}},
	}
	nopLogger := zap.NewNop().Sugar()
	svc := NewExchangeService(clients, 1*time.Second, nopLogger)
	q := svc.BestQuote(context.Background(), domain.CurrencyPair{})
	if q.Provider != "b" || q.Rate != 3.0 {
		t.Errorf("expected best b with rate 3.0, got %v %.2f", q.Provider, q.Rate)
	}
}

func TestBestQuote_AllFail(t *testing.T) {
	clients := []client.RateClient{
		&mockClient{quote: domain.Quote{Provider: "a", Err: errors.New("fail1")}},
		&mockClient{quote: domain.Quote{Provider: "b", Err: errors.New("fail2")}},
	}
	nopLogger := zap.NewNop().Sugar()
	svc := NewExchangeService(clients, 1*time.Second, nopLogger)
	q := svc.BestQuote(context.Background(), domain.CurrencyPair{})
	if q.Err == nil {
		t.Error("expected error when all providers fail")
	}
}
