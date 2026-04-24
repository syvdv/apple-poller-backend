package models

type ITunesFeed struct {
	Feed struct {
		Entry []ReviewEntry `json:"entry"`
	} `json:"feed"`
}

type ReviewEntry struct {
	ID struct {
		Label string `json:"label"`
	} `json:"id"`
	Author struct {
		Name struct {
			Label string `json:"label"`
		} `json:"name"`
	} `json:"author"`
	Rating struct {
		Label string `json:"label"`
	} `json:"im:rating"`
	Content struct {
		Label string `json:"label"`
	} `json:"content"`
	Updated struct {
		Label string `json:"label"`
	} `json:"updated"`
}
