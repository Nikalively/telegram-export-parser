package telegram_export_parser

import (
	"time"
)

// Event represents a single message from Telegram export.
type Event struct {
	ID       int64     `json:"id"`
	FromID   string    `json:"from_id"` // Unique ID from JSON, used instead of username
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
	ChatID   int64     `json:"chat_id"` // ID of the chat
	Entities []Entity  `json:"entities,omitempty"`
}

// Entity represents text entities like mentions.
type Entity struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// MergedEvents holds a sorted list of unique events.
type MergedEvents struct {
	Events []Event `json:"events"`
}
