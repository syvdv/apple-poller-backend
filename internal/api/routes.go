package api

import (
	"net/http"

	"applepoller/internal/config"
)

func RegisterRoutes(cfg *config.AppConfig) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/apps", GetAppIDsHandler)
	mux.HandleFunc("/api/reviews", GetReviewsHandler(cfg.GetMaxAge()))

	return EnableCORS(cfg.CorsAllowedOrigin, mux)
}
