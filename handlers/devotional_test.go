package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lwnra-devo-api/database"
	"lwnra-devo-api/facebook"
)

func TestParseDevotional(t *testing.T) {
	// Create test handler
	db, _ := database.New(":memory:")
	fbClient := facebook.New("")
	handler := NewDevotionalHandler(db, fbClient)

	// Test request
	requestBody := map[string]string{
		"message": `DAILY DEVOTIONAL
Read Matthew 6:16-18
August 2, 2025
Matthew 6:16-18 NIV
16 When you fast, do not look somber as the hypocrites do, for they disfigure their faces to show men they are fasting.
REFLECTION QUESTIONS
What spiritual habit do you do partly for others to notice?`,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/devotionals/parse", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	handler.ParseDevotional(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	// Check if devotional was parsed correctly
	devotionalData := response.Data.(map[string]interface{})
	if devotionalData["reading"] != "Matthew 6:16-18" {
		t.Errorf("Expected reading 'Matthew 6:16-18', got %v", devotionalData["reading"])
	}

	if devotionalData["version"] != "NIV" {
		t.Errorf("Expected version 'NIV', got %v", devotionalData["version"])
	}

	if devotionalData["date"] != "August 2, 2025" {
		t.Errorf("Expected date 'August 2, 2025', got %v", devotionalData["date"])
	}
}
