package models

type Review struct {
	AppID   string `json:"app_id"`
	ID      string `json:"id"`
	Author  string `json:"author"`
	Score   string `json:"score"`
	Content string `json:"content"`
	Time    string `json:"time"`
}
