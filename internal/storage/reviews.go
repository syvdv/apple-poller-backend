package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"applepoller/internal/models"
)

func GetRecentReviews(appID string, maxAge time.Duration) []models.Review {
	// Initialize as an empty slice so it returns [] instead of null in JSON.
	results := []models.Review{}

	file, err := os.Open("reviews.jsonl")
	if err != nil {
		return results // If file doesn't exist yet, return empty list.
	}
	defer file.Close()

	cutoffTime := time.Now().Add(-maxAge)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rev models.Review
		if err := json.Unmarshal(scanner.Bytes(), &rev); err == nil {
			if rev.AppID == appID {
				revTime, err := time.Parse(time.RFC3339, rev.Time)

				if err == nil && revTime.After(cutoffTime) {
					results = append(results, rev)
				}
			}
		}
	}

	return results
}

func LoadLastReviewIDs() map[string]string {
	latestIDs := make(map[string]string)
	latestTimes := make(map[string]string)

	file, err := os.Open("reviews.jsonl")
	if err != nil {
		return latestIDs
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var rev models.Review

		if err := json.Unmarshal(scanner.Bytes(), &rev); err == nil {
			// Apple's time strings are ISO-8601, they sort alphabetically.
			if rev.Time >= latestTimes[rev.AppID] {
				latestTimes[rev.AppID] = rev.Time
				latestIDs[rev.AppID] = rev.ID
			}
		}
	}

	return latestIDs
}

func SaveReviews(reviews []models.Review) error {
	if len(reviews) == 0 {
		return nil
	}

	file, err := os.OpenFile("reviews.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, rev := range reviews {
		data, _ := json.Marshal(rev)
		file.Write(data)
		file.WriteString("\n")
	}

	return nil
}
