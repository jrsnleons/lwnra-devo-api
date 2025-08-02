package main

import (
	"fmt"
	"os"
	"time"

	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
	"lwnra-devo-api/parser"
)

func main() {
	// Get Facebook access token from environment
	accessToken := os.Getenv("FB_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("Please set FB_ACCESS_TOKEN as an environment variable")
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New("devotionals.db")
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize Facebook client
	fbClient := facebook.New(accessToken)

	// Fetch recent posts
	posts, err := fbClient.GetRecentPosts()
	if err != nil {
		fmt.Printf("Failed to fetch Facebook posts: %v\n", err)
		os.Exit(1)
	}

	// Filter for devotional posts from today/yesterday
	devotionalPosts := facebook.FilterDevotionalPosts(posts)

	// Process each devotional post
	count := 0
	for _, post := range devotionalPosts {
		devo := parser.ParseDevotional(post.Message)

		// Debug: Print what the parser found
		fmt.Printf("Debug - Parsed date: '%s'\n", devo.Date)

		// Use post date if devotional date is empty
		if devo.Date == "" {
			// Extract date from Facebook post's created_time and format it properly
			postDate := extractAndFormatPostDate(post.CreatedTime)
			if postDate != "" {
				devo.Date = postDate
				fmt.Printf("Debug - Using post date: '%s'\n", postDate)
			} else {
				// Only use "today" as final fallback
				devo.Date = "today"
				fmt.Printf("Debug - Using fallback: 'today'\n")
			}
		}

		// Save to database
		if err := db.SaveDevotional(devo); err != nil {
			fmt.Printf("Failed to save devotional '%s': %v\n", devo.Title, err)
		} else {
			fmt.Printf("[%s] %s — Saved to DB\n", devo.Date, devo.Title)
			count++
		}
	}

	if count == 0 {
		fmt.Println("No new DAILY DEVOTIONAL found for today or yesterday.")
	} else {
		fmt.Printf("Successfully processed %d devotional(s)\n", count)
	}
}

// extractAndFormatPostDate converts Facebook's created_time to a readable date format
func extractAndFormatPostDate(createdTime string) string {
	t, err := time.Parse("2006-01-02T15:04:05-0700", createdTime)
	if err != nil {
		return ""
	}

	// Convert to "August 2, 2025" format
	return t.Format("January 2, 2006")
}
