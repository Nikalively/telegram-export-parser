package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Nikalively/telegram-export-parser"
)

// ParseSingleFile parses a single result.json file and returns events.
func ParseSingleFile(filePath string) ([]telegram_export_parser.Event, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rawMessages []RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var events []telegram_export_parser.Event
	for _, msg := range rawMessages {
		event := telegram_export_parser.Event{
			ID:     msg.ID,
			FromID: msg.FromID,
			Text:   msg.Text,
			Date:   msg.Date,
			ChatID: msg.Chat.ID,
		}
		if msg.Entities != nil {
			event.Entities = msg.Entities
		}
		events = append(events, event)
	}
	return events, nil
}

// MergeFromFolders merges events from multiple export folders, sorts by date, and deduplicates by ID.
func MergeFromFolders(folderPaths []string) (*telegram_export_parser.MergedEvents, error) {
	eventMap := make(map[int64]telegram_export_parser.Event)
	for _, folder := range folderPaths {
		files, err := filepath.Glob(filepath.Join(folder, "result.json"))
		if err != nil {
			return nil, fmt.Errorf("failed to glob files in %s: %w", folder, err)
		}
		for _, file := range files {
			events, err := ParseSingleFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", file, err)
			}
			for _, event := range events {
				eventMap[event.ID] = event // Deduplicate by ID
			}
		}
	}

	events := make([]telegram_export_parser.Event, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})

	return &telegram_export_parser.MergedEvents{Events: events}, nil
}

// RawMessage is the raw structure from Telegram JSON export.
type RawMessage struct {
	ID       int64                           `json:"id"`
	FromID   string                          `json:"from_id"`
	Text     string                          `json:"text"`
	Date     time.Time                       `json:"date"`
	Chat     RawChat                         `json:"chat"`
	Entities []telegram_export_parser.Entity `json:"entities,omitempty"`
}

// RawChat represents the chat in raw JSON.
type RawChat struct {
	ID int64 `json:"id"`
}
