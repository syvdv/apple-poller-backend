package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"applepoller/internal/storage"
)

func GetAppIDsHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load the current state of last review IDs from storage.
	state := storage.LoadLastReviewIDs()
	appIDs := make([]string, 0, len(state))
	for id := range state {
		appIDs = append(appIDs, id)
	}

	if appIDs == nil {
		appIDs = []string{}
	}

	// Sort app IDs alphabetically for consistent output.
	sort.Strings(appIDs)

	// Return the app IDs as JSON.
	w.Header().Set("Content-Type", "application/json")
	response := map[string][]string{
		"app_ids": appIDs,
	}

	json.NewEncoder(w).Encode(response)
}

func GetReviewsHandler(maxAge time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests.
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the app ID from the query parameters.
		appID := r.URL.Query().Get("id")
		if appID == "" {
			http.Error(w, "Missing 'id' parameter in URL", http.StatusBadRequest)
			return
		}

		// Check if the app ID exists in our storage. If not, return a 404 error.
		state := storage.LoadLastReviewIDs()
		if _, exists := state[appID]; !exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unknown app ID: " + appID,
			})
			return
		}

		// Get recent reviews for the app ID, filtered by maxAge.
		reviews := storage.GetRecentReviews(appID, maxAge)

		// Sort reviews by time in descending order (newest first).
		sort.Slice(reviews, func(i, j int) bool {
			return reviews[i].Time > reviews[j].Time
		})

		// Return the reviews as JSON.
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"app_id":  appID,
			"reviews": reviews,
		}

		json.NewEncoder(w).Encode(response)
	}
}
