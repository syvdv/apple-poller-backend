package config

import (
	"encoding/json"
	"os"
	"time"
)

type AppConfig struct {
	CorsAllowedOrigin string   `json:"cors_allowed_origin"`
	APIPort           string   `json:"api_port"`
	APIMaxReviewsAge  string   `json:"api_max_reviews_age"`
	URLTemplate       string   `json:"url_template"`
	PollInterval      string   `json:"poll_interval"`
	PaginationLimit   int      `json:"pagination_limit"` // Apple limits to 10 pages.
	AppIDs            []string `json:"app_ids"`
}

func (c *AppConfig) GetMaxAge() time.Duration {
	duration, err := time.ParseDuration(c.APIMaxReviewsAge)
	if err != nil {
		// Fallback to 48 hours if the string is invalid.
		return 48 * time.Hour
	}
	return duration
}

func (c *AppConfig) GetDuration() time.Duration {
	duration, err := time.ParseDuration(c.PollInterval)
	if err != nil {
		// Fallback to 10 minutes if the string is invalid.
		return 10 * time.Minute
	}
	return duration
}

func Load() (*AppConfig, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
