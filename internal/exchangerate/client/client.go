package client

import (
	"context"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type RateClient interface {
	FetchRate(ctx context.Context, pair domain.CurrencyPair) domain.Quote
}
