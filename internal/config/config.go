package config

import (
	"os"
	"strconv"
)

type Config struct {
	KafkaBrokers string
	Port         string
	AuthToken    string
	RateLimit    int
}

func Load() (*Config, error) {
	rateLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if rateLimit == 0 {
		rateLimit = 100 // default rate limit
	}

	return &Config{
		KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
		Port:         os.Getenv("PORT"),
		AuthToken:    os.Getenv("AUTH_TOKEN"),
		RateLimit:    rateLimit,
	}, nil
}
