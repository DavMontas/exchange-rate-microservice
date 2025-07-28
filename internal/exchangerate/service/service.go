package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/client"
	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type ExchangeService struct {
	Clients []client.RateClient
	Timeout time.Duration
	Log     *zap.SugaredLogger
}

func NewExchangeService(clients []client.RateClient, timeout time.Duration, log *zap.SugaredLogger) *ExchangeService {
	return &ExchangeService{Clients: clients, Timeout: timeout, Log: log}
}

func (s *ExchangeService) BestQuote(ctx context.Context, pair domain.CurrencyPair) domain.Quote {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	ch := make(chan domain.Quote, len(s.Clients))
	for _, rc := range s.Clients {
		go func(rc client.RateClient) {
			ch <- rc.FetchRate(ctx, pair)
		}(rc)
	}

	var best domain.Quote
	var found bool
	results := 0
	for results < len(s.Clients) {
		select {
		case q := <-ch:
			results++
			if q.Err != nil {
				s.Log.Warnw("provider error", "provider", q.Provider, "error", q.Err)
				continue
			}
			s.Log.Infow("provider result", "provider", q.Provider, "rate", q.Rate)
			if !found || q.Rate > best.Rate {
				best = q
				found = true
			}
		case <-ctx.Done():
			s.Log.Warnw("timeout reached")
			results = len(s.Clients)
		}
	}

	if !found {
		err := fmt.Errorf("no valid rate found for %s to %s", pair.From, pair.To)
		s.Log.Errorw("all providers failed", "error", err)
		return domain.Quote{Err: err}
	}
	return best
}
