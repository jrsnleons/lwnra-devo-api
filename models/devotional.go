package models

// Devotional represents a daily devotional entry
type Devotional struct {
	Date         string   `json:"date"`          // e.g. "August 2, 2025"
	Reading      string   `json:"reading"`       // "Matthew 6:16-18"
	Version      string   `json:"version"`       // Bible version like "NIV", "ESV", "NASB"
	Passage      string   `json:"passage"`       // passage text
	ReflectionQs []string `json:"reflection_qs"` // questions
	Title        string   `json:"title"`         // devo title
	Author       string   `json:"author"`        // author
	Body         string   `json:"body"`          // main devo body
	Prayer       string   `json:"prayer"`        // prayer
}
