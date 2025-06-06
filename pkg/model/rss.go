package model

import (
	"strings"
)

type RSS struct {
	Channel RSSChannel `xml:"channel"`
	Type    string
}

type RSSChannel struct {
	Title       string           `xml:"title"`
	Link        string           `xml:"link"`
	Description string           `xml:"description"`
	Items       []RSSChannelItem `xml:"item"`
}

type RSSChannelItem struct {
	Title    string `xml:"title"`
	Link     string `xml:"link"`
	Content  string `xml:"encoded" xml:"content:encoded"`
	Category string `xml:"category"`
	PubDate  string `xml:"pubDate"`
	GUID     string `xml:"guid"`
	DCDate   string `xml:"date" xml:"dc:date"`
}

type RSSItem struct {
	Title   string   `json:"title"`
	Link    string   `json:"link"`
	Content string   `json:"content"`
	Type    []string `json:"type"`
	GUID    string   `json:"guid"`
}

func (r RSSItem) ToMessage() Message {
	content := strings.ReplaceAll(r.Content, "<p>", "")
	content = strings.ReplaceAll(content, "</p>", "")

	return Message{
		Type:    strings.Join(r.Type, ","),
		Title:   r.Title,
		Content: content,
		Link:    r.Link,
	}
}
