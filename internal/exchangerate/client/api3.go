package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
)

type xmlEnvelope struct {
	Cube struct {
		Cube []struct {
			Currency string  `xml:"currency,attr"`
			Rate     float64 `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}

type API3Client struct {
	URL    string
	Client *http.Client
}

func NewAPI3Client(url string) *API3Client {
	return &API3Client{
		URL:    url,
		Client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *API3Client) FetchRate(ctx context.Context, pair domain.CurrencyPair) domain.Quote {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}
	defer resp.Body.Close()

	var env xmlEnvelope
	if err := xml.NewDecoder(resp.Body).Decode(&env); err != nil {
		return domain.Quote{Provider: "api3", Err: err}
	}

	from := strings.ToUpper(pair.From)
	to := strings.ToUpper(pair.To)

	if from == "EUR" {
		for _, cube := range env.Cube.Cube {
			if cube.Currency == to {
				return domain.Quote{Provider: "api3", Rate: cube.Rate * pair.Amount}
			}
		}
		return domain.Quote{Provider: "api3", Err: fmt.Errorf("currency %s not found", to)}
	}

	for _, cube := range env.Cube.Cube {
		if cube.Currency == from {
			return domain.Quote{Provider: "api3", Rate: (1.0 / cube.Rate) * pair.Amount}
		}
	}
	return domain.Quote{Provider: "api3", Err: fmt.Errorf("currency %s not found", from)}
}
