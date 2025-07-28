package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type API3Client struct {
	URL    string
	Client *http.Client
}

func NewAPI3Client(url string) *API3Client {
	return &API3Client{
		URL:    url,
		Client: &http.Client{Timeout: 1 * time.Second},
	}
}

func (c *API3Client) FetchRate(ctx context.Context, pair domain.CurrencyPair) domain.Quote {
	endpoint := fmt.Sprintf(c.URL, pair.From, pair.To)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}
	defer resp.Body.Close()

	var out struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}

	ratePerUnit, ok := out.Rates[pair.To]
	if !ok {
		return domain.Quote{Provider: "api3", Err: fmt.Errorf("currency %s not found", pair.To)}
	}

	return domain.Quote{Provider: "api3", Rate: ratePerUnit * pair.Amount}
}


