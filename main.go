package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Puzzle struct {
	ID        int    `json:"puzzle_id"`
	Category  string `json:"category"`
	URL       string `json:"url"`
	LichessID string `json:"lichess_id"`
}

// Global variable to store puzzles loaded from the database
var puzzles []Puzzle

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Get a database connection
func getDBConnection() (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_PUBLIC_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_PUBLIC_URL is not set in the environment")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Query all puzzles from the database on initial startup
func loadPuzzlesFromDB() {
	conn, err := getDBConnection()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer conn.Close(context.Background())

	query := `SELECT puzzle_id, category, url FROM puzzles;`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to query puzzles:", err)
	}
	defer rows.Close()

	var puzzleList []Puzzle

	// ADDED: Track unique IDs during initial load
	idMap := make(map[string]bool)
	duplicates := make([]string, 0)

	for rows.Next() {
		var p Puzzle
		if err := rows.Scan(&p.ID, &p.Category, &p.URL); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Extract Lichess ID from URL
		parts := strings.Split(p.URL, "/")
		if len(parts) > 0 {
			p.LichessID = parts[len(parts)-1]
			// ADDED: Check for duplicates during load
			if idMap[p.LichessID] {
				duplicates = append(duplicates, p.LichessID)
			}
			idMap[p.LichessID] = true
		} else {
			log.Printf("Invalid URL format for puzzle ID %d: %s", p.ID, p.URL)
			continue
		}

		puzzleList = append(puzzleList, p)
	}

	// ADDED: Detailed logging of puzzle load results
	log.Printf("Database load summary:")
	log.Printf("- Total puzzles loaded: %d", len(puzzleList))
	log.Printf("- Unique Lichess IDs: %d", len(idMap))
	if len(duplicates) > 0 {
		log.Printf("WARNING: Found %d duplicate lichess_ids in database: %v", len(duplicates), duplicates)
	}

	// Log category distribution
	categoryCount := make(map[string]int)
	for _, p := range puzzleList {
		categoryCount[p.Category]++
	}
	log.Println("Category distribution:")
	for category, count := range categoryCount {
		log.Printf("- %s: %d puzzles", category, count)
	}

	// Store the puzzles in the global variable without shuffling
	puzzles = puzzleList
}

// the API endpoint to retrieve and shuffle the 100 puzzles
func getPuzzles(w http.ResponseWriter, r *http.Request) {
	requestID := rand.Int()
	log.Printf("[Request %d] Generating puzzle set", requestID)

	// Create a copy of indices and shuffle them
	indices := make([]int, len(puzzles))
	for i := range indices {
		indices[i] = i
	}
	for i := len(indices) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	// Take the first 100 (or less if fewer puzzles exist)
	result := make([]Puzzle, min(100, len(puzzles)))
	for i := range result {
		result[i] = puzzles[indices[i]]
	}

	// Validate response set
	idMap := make(map[string]bool)
	duplicates := make([]string, 0)
	categoryCount := make(map[string]int)

	for _, p := range result {
		if idMap[p.LichessID] {
			duplicates = append(duplicates, p.LichessID)
		}
		idMap[p.LichessID] = true
		categoryCount[p.Category]++
	}

	// ADDED: Detailed response logging
	log.Printf("[Request %d] Response summary:", requestID)
	log.Printf("- Puzzles in response: %d", len(result))
	log.Printf("- Unique puzzles: %d", len(idMap))
	if len(duplicates) > 0 {
		log.Printf("- WARNING: Found duplicates: %v", duplicates)
	}
	log.Println("- Category distribution in response:")
	for category, count := range categoryCount {
		log.Printf("  %s: %d puzzles", category, count)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	loadPuzzlesFromDB() // Load puzzles once at startup

	mux := http.NewServeMux()
	mux.HandleFunc("/api/puzzles", getPuzzles)

	log.Println("Puzzles API Server is running on port 8081...")
	http.ListenAndServe(":8081", withCORS(mux))
}
