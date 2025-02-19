// main_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test setup helper
func setupTestPuzzles() {
	// Initialize test puzzle data
	puzzles = []Puzzle{
		{ID: 1, Category: "Pins", URL: "https://lichess.org/training/abc123", LichessID: "abc123"},
		{ID: 2, Category: "Forks", URL: "https://lichess.org/training/def456", LichessID: "def456"},
		{ID: 3, Category: "Skewers", URL: "https://lichess.org/training/ghi789", LichessID: "ghi789"},
	}
}

func TestGetPuzzles(t *testing.T) {
	setupTestPuzzles()

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/api/puzzles", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPuzzles)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Parse the response body
	var response []Puzzle
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	// Verify response length
	expectedLength := min(100, len(puzzles))
	if len(response) != expectedLength {
		t.Errorf("handler returned wrong number of puzzles: got %v want %v",
			len(response), expectedLength)
	}

	// Verify that all returned puzzles are valid
	for _, p := range response {
		if p.ID == 0 || p.Category == "" || p.URL == "" || p.LichessID == "" {
			t.Errorf("handler returned invalid puzzle: %+v", p)
		}
	}
}

func TestWithCORS(t *testing.T) {
	setupTestPuzzles()

	// Create a test server with CORS middleware
	handler := withCORS(http.HandlerFunc(getPuzzles))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request to the test server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check CORS headers
	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, expected := range expectedHeaders {
		if got := resp.Header.Get(header); got != expected {
			t.Errorf("handler returned wrong %s header: got %v want %v",
				header, got, expected)
		}
	}
}
