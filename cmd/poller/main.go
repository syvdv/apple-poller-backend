package main

import (
	"log"
	"net/http"
	"time"

	"applepoller/internal/api"
	"applepoller/internal/config"
	"applepoller/internal/fetcher"
	"applepoller/internal/storage"
)

func main() {
	log.Println("Starting Apple Review Poller.")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config.json: %v", err)
	}

	interval := cfg.GetDuration()

	log.Printf("Config loaded.")
	log.Printf("Tracking %d apps every %v", len(cfg.AppIDs), interval)

	// Start the polling loop in a separate goroutine.
	go func() {
		for {
			currentState := storage.LoadLastReviewIDs()

			for _, appID := range cfg.AppIDs {
				fetcher.PollApp(appID, cfg.URLTemplate, cfg.PaginationLimit, currentState)
			}

			log.Printf("Polling goes to sleep for %v.", interval)
			time.Sleep(interval)
		}
	}()

	// Register API routes and start the server.
	router := api.RegisterRoutes(cfg)

	log.Printf("Starting API server on http://localhost%s.", cfg.APIPort)
	http.ListenAndServe(":8080", router)
}
