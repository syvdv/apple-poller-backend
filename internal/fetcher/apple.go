package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"applepoller/internal/models"
	"applepoller/internal/storage"
)

func PollApp(appID string, urlTemplate string, PaginationLimit int, state map[string]string) {
	log.Printf("Fetching reviews for App ID %s.", appID)

	lastSeenID := state[appID]
	var newReviews []models.Review

	for page := 1; page <= PaginationLimit; page++ {
		url := fmt.Sprintf(urlTemplate, appID, page)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching page %d for %s: %v", page, appID, err)
			break
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error reading response body on page %d: %v", page, err)
			break
		}

		var feedData models.ITunesFeed
		if err := json.Unmarshal(body, &feedData); err != nil {
			log.Printf("Error parsing JSON on page %d: %v", page, err)
			break
		}

		if len(feedData.Feed.Entry) == 0 {
			break
		}

		caughtUp := false

		for _, entry := range feedData.Feed.Entry {
			// If we hit a review we've already seen, we are fully caught up.
			if entry.ID.Label == lastSeenID {
				caughtUp = true
				break
			}

			clean := models.Review{
				AppID:   appID,
				ID:      entry.ID.Label,
				Author:  entry.Author.Name.Label,
				Score:   entry.Rating.Label,
				Content: entry.Content.Label,
				Time:    entry.Updated.Label,
			}
			newReviews = append(newReviews, clean)
		}

		if caughtUp {
			break
		}
	}

	if len(newReviews) > 0 {
		log.Printf("Found %d new reviews for App ID %s.", len(newReviews), appID)
		storage.SaveReviews(newReviews)
	} else {
		log.Printf("No new reviews for App ID %s.", appID)
	}
}
