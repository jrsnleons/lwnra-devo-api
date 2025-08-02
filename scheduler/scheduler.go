package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
	"lwnra-devo-api/parser"
)

// Scheduler handles automated tasks
type Scheduler struct {
	cron *cron.Cron
	db   *database.DB
	fb   *facebook.Client
}

// New creates a new scheduler instance
func New(db *database.DB, fb *facebook.Client) *Scheduler {
	// Use Philippine timezone (UTC+8)
	philippineLocation, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		log.Printf("Failed to load Philippine timezone, using UTC: %v", err)
		philippineLocation = time.UTC
	}

	c := cron.New(cron.WithLocation(philippineLocation))

	return &Scheduler{
		cron: c,
		db:   db,
		fb:   fb,
	}
}

// Start begins the scheduled tasks
func (s *Scheduler) Start() {
	// Schedule devotional sync at 4:45 AM Philippine time every day
	// Cron format: "45 4 * * *" = 45 minutes, 4 hours, every day of month, every month, every day of week
	_, err := s.cron.AddFunc("45 4 * * *", s.syncDevotionals)
	if err != nil {
		log.Printf("Failed to schedule devotional sync: %v", err)
		return
	}

	// Optional: Also run a backup sync at 5:15 AM in case the first one fails
	_, err = s.cron.AddFunc("15 5 * * *", s.syncDevotionals)
	if err != nil {
		log.Printf("Failed to schedule backup devotional sync: %v", err)
	}

	s.cron.Start()
	log.Println("Scheduler started - devotionals will sync daily at 4:45 AM Philippine time")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// syncDevotionals performs the automated devotional sync
func (s *Scheduler) syncDevotionals() {
	log.Println("Starting scheduled devotional sync...")

	// Get recent posts from Facebook
	posts, err := s.fb.GetRecentPosts()
	if err != nil {
		log.Printf("Failed to fetch Facebook posts during scheduled sync: %v", err)
		return
	}

	// Filter for devotional posts from today and yesterday
	devotionalPosts := facebook.FilterDevotionalPosts(posts)
	
	syncCount := 0
	for _, post := range devotionalPosts {
		// Parse the devotional content
		devotional := parser.ParseDevotional(post.Message)

		// Save to database
		err = s.db.SaveDevotional(devotional)
		if err != nil {
			log.Printf("Failed to save devotional during scheduled sync: %v", err)
			continue
		}

		syncCount++
	}

	if syncCount > 0 {
		log.Printf("Scheduled sync completed successfully - %d devotionals synced", syncCount)
	} else {
		log.Println("Scheduled sync completed - no new devotionals found")
	}
}

// GetNextRun returns the next scheduled run time
func (s *Scheduler) GetNextRun() time.Time {
	entries := s.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}

// IsRunning returns whether the scheduler is currently running
func (s *Scheduler) IsRunning() bool {
	return len(s.cron.Entries()) > 0
}
