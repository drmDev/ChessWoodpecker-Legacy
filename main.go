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

var puzzles []Puzzle

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
	seenIDs := make(map[string]bool)
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
			// Check for duplicates
			if seenIDs[p.LichessID] {
				duplicates = append(duplicates, p.LichessID)
			}
			seenIDs[p.LichessID] = true
		} else {
			log.Printf("Invalid URL format for puzzle ID %d: %s", p.ID, p.URL)
			continue
		}

		puzzleList = append(puzzleList, p)
	}

	// Validate puzzle count
	if len(puzzleList) < 200 {
		log.Printf("⚠️ Warning: Expected 200 puzzles but only loaded %d", len(puzzleList))
	} else {
		log.Printf("✅ Successfully loaded %d puzzles", len(puzzleList))
	}

	// Report any duplicates
	if len(duplicates) > 0 {
		log.Printf("⚠️ Warning: Found duplicate Lichess IDs: %v", duplicates)
	}

	puzzles = puzzleList
}

func getPuzzles(w http.ResponseWriter, r *http.Request) {
	// Create and shuffle indices
	indices := make([]int, len(puzzles))
	for i := range indices {
		indices[i] = i
	}
	for i := len(indices) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	result := make([]Puzzle, len(puzzles))
	for i := range result {
		result[i] = puzzles[indices[i]]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	loadPuzzlesFromDB() // Load puzzles once at startup

	mux := http.NewServeMux()
	mux.HandleFunc("/api/puzzles", getPuzzles)

	log.Printf("✅ Server initialized with %d puzzles", len(puzzles))
	log.Println("Puzzles API Server is running on port 8081...")
	http.ListenAndServe(":8081", withCORS(mux))
}
