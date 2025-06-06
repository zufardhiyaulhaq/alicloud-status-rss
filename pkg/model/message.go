package model

type Message struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Link    string `json:"link"`
}
