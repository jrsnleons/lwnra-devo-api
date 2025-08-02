package facebook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lwnra-devo-api/models"
)

// Client handles Facebook API interactions
type Client struct {
	accessToken string
}

// New creates a new Facebook client
func New(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
	}
}

// GetRecentPosts fetches recent posts from Facebook API
func (c *Client) GetRecentPosts() ([]models.FBPost, error) {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v23.0/me?fields=id,name,posts{message,created_time}&access_token=%s",
		c.accessToken,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make Facebook API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Facebook API error (status %d): %s", resp.StatusCode, string(body))
	}

	var fb models.FBMeResponse
	if err := json.NewDecoder(resp.Body).Decode(&fb); err != nil {
		return nil, fmt.Errorf("failed to decode Facebook API response: %w", err)
	}

	return fb.Posts.Data, nil
}

// FilterDevotionalPosts filters posts to only include daily devotionals from today or yesterday
func FilterDevotionalPosts(posts []models.FBPost) []models.FBPost {
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	var filtered []models.FBPost
	for _, post := range posts {
		if !isDevotionalPost(post.Message) {
			continue
		}

		postDate := extractPostDate(post.CreatedTime)
		if postDate == today || postDate == yesterday {
			filtered = append(filtered, post)
		}
	}

	return filtered
}

// isDevotionalPost checks if a post is a daily devotional
func isDevotionalPost(message string) bool {
	return len(message) > 0 && (message[:min(16, len(message))] == "DAILY DEVOTIONAL")
}

// extractPostDate extracts the date from Facebook's created_time format
func extractPostDate(createdTime string) string {
	t, err := time.Parse("2006-01-02T15:04:05-0700", createdTime)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
