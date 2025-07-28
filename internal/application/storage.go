package application

import "github.com/davmontas/exchange-rate-offers/internal/exchangerate/client"

func RegisterAPIs(cfg Config) []client.RateClient {
	return []client.RateClient{
		client.NewAPI1Client(cfg.Storage.API1.URL),
		client.NewAPI2Client(cfg.Storage.API2.URL),
		client.NewAPI3Client(cfg.Storage.API3.URL),
	}
}
