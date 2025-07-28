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

// XML sample:
// <Envelope>
//   <Cube>
//     <Cube time="2025-07-27">
//       <Cube currency="USD" rate="1.089" />
//       <Cube currency="JPY" rate="168.3" />
//     </Cube>
//   </Cube>
// </Envelope>

type xmlEnvelope struct {
	Cube struct {
		Times []struct {
			Cubes []struct {
				Currency string  `xml:"currency,attr"`
				Rate     float64 `xml:"rate,attr"`
			} `xml:"Cube"`
		} `xml:"Cube"`
	} `xml:"Cube"`
}

type API2Client struct {
	URL    string
	Client *http.Client
}

func NewAPI2Client(url string) *API2Client {
	return &API2Client{
		URL:    url,
		Client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (c *API2Client) FetchRate(ctx context.Context, pair domain.CurrencyPair) domain.Quote {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return domain.Quote{Provider: "api2", Err: err}
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return domain.Quote{Provider: "api2", Err: err}
	}
	defer resp.Body.Close()

	var env xmlEnvelope
	if err := xml.NewDecoder(resp.Body).Decode(&env); err != nil {
		return domain.Quote{Provider: "api2", Err: err}
	}

	from := strings.ToUpper(pair.From)
	to := strings.ToUpper(pair.To)

	// El XML asume que todas las tasas están en base EUR
	if from == "EUR" {
		for _, timeCube := range env.Cube.Times {
			for _, cube := range timeCube.Cubes {
				if strings.ToUpper(cube.Currency) == to {
					return domain.Quote{Provider: "api2", Rate: cube.Rate * pair.Amount}
				}
			}
		}
		return domain.Quote{Provider: "api2", Err: fmt.Errorf("currency %s not found", to)}
	}

	// Si el 'from' no es EUR, buscar su tasa para convertir a EUR
	for _, timeCube := range env.Cube.Times {
		for _, cube := range timeCube.Cubes {
			if strings.ToUpper(cube.Currency) == from {
				rate := 1.0 / cube.Rate
				return domain.Quote{Provider: "api2", Rate: rate * pair.Amount}
			}
		}
	}
	return domain.Quote{Provider: "api2", Err: fmt.Errorf("currency %s not found", from)}
}
