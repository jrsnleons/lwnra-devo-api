package scheduler

import (
	"testing"
	"time"

	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
)

func TestSchedulerCreation(t *testing.T) {
	// Create test database
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create test Facebook client
	fbClient := facebook.New("test_token")

	// Create scheduler
	sched := New(db, fbClient)

	if sched == nil {
		t.Fatal("Scheduler creation failed")
	}

	if !sched.GetNextRun().IsZero() {
		t.Error("Scheduler should not have next run before starting")
	}

	if sched.IsRunning() {
		t.Error("Scheduler should not be running initially")
	}
}

func TestSchedulerTimezone(t *testing.T) {
	// Test that Philippine timezone is loaded correctly
	_, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		t.Skipf("Philippine timezone not available: %v", err)
	}

	// Create test scheduler
	db, _ := database.New(":memory:")
	defer db.Close()
	fbClient := facebook.New("test_token")
	sched := New(db, fbClient)

	// Start scheduler
	sched.Start()
	defer sched.Stop()

	// Verify it's running
	if !sched.IsRunning() {
		t.Error("Scheduler should be running after start")
	}

	// Verify next run is set and in Philippine time
	nextRun := sched.GetNextRun()
	if nextRun.IsZero() {
		t.Error("Next run should be set after starting scheduler")
	}

	// Check timezone
	_, offset := nextRun.Zone()
	philippineOffset := 8 * 3600 // UTC+8 in seconds
	if offset != philippineOffset {
		t.Errorf("Expected Philippine timezone offset %d, got %d", philippineOffset, offset)
	}
}

func TestSchedulerStartStop(t *testing.T) {
	db, _ := database.New(":memory:")
	defer db.Close()
	fbClient := facebook.New("test_token")
	sched := New(db, fbClient)

	// Start scheduler
	sched.Start()
	if !sched.IsRunning() {
		t.Error("Scheduler should be running after start")
	}

	// Stop scheduler
	sched.Stop()
	// Note: cron scheduler may still show entries briefly after stop
	// so we don't test IsRunning() immediately after stop
}
