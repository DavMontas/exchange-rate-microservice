package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type API1Client struct {
	URL    string
	Client *http.Client
}

func NewAPI1Client(url string) *API1Client {
	return &API1Client{URL: url, Client: &http.Client{Timeout: 2 * time.Second}}
}

func (c *API1Client) FetchRate(ctx context.Context, pair domain.CurrencyPair) domain.Quote {
	endpoint := fmt.Sprintf(c.URL, strings.ToLower(pair.From))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.Quote{Provider: "api1", Err: err}
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return domain.Quote{Provider: "api1", Err: err}
	}
	defer resp.Body.Close()

	rates := map[string]struct {
		Rate float64 `json:"rate"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return domain.Quote{Provider: "api1", Err: err}
	}
	entry, ok := rates[strings.ToLower(pair.To)]
	if !ok {
		return domain.Quote{Provider: "api1", Err: fmt.Errorf("currency %s not found", pair.To)}
	}
	return domain.Quote{Provider: "api1", Rate: entry.Rate * pair.Amount}
}
