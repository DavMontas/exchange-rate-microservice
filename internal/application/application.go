package application

import "time"

type ServerConfig struct {
	ListenAddr string
	Mode       string
}

type APIConfig struct {
	URL string
}

type ServiceConfig struct {
	Timeout time.Duration
}

type Config struct {
	Server  ServerConfig
	Storage Storage
	Service ServiceConfig
}

type Storage struct {
	API1 APIConfig
	API2 APIConfig
	API3 APIConfig
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			ListenAddr: ":8080",
			Mode:       "release",
		},
		Storage: Storage{
			API1: APIConfig{
				URL: "http://www.floatrates.com/daily/%s.json",
			},
			API2: APIConfig{
				URL: "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
			},
			API3: APIConfig{
				URL: "https://api.frankfurter.app/latest?from=%s&to=%s",

			},
		},
		Service: ServiceConfig{
			Timeout: 3 * time.Second,
		},
	}
}
