package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// Puzzle struct matching your database schema
type Puzzle struct {
	ID       int    `json:"puzzle_id"`
	Category string `json:"category"`
	URL      string `json:"url"`
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

// Query all puzzles from the database
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
	for rows.Next() {
		var p Puzzle
		if err := rows.Scan(&p.ID, &p.Category, &p.URL); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		puzzleList = append(puzzleList, p)
	}

	// Store the puzzles in the global variable without shuffling
	puzzles = puzzleList
	log.Println("Loaded", len(puzzles), "puzzles from database")
}

// API handler to return randomized puzzles
func getPuzzles(w http.ResponseWriter, r *http.Request) {
	// Shuffle puzzles randomly on each request
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(puzzles), func(i, j int) { puzzles[i], puzzles[j] = puzzles[j], puzzles[i] })

	// Limit to the first 100 puzzles
	if len(puzzles) > 100 {
		puzzles = puzzles[:100]
	}

	log.Println("Serving randomized puzzles")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(puzzles)
}

func main() {
	loadPuzzlesFromDB() // Load puzzles once at startup

	mux := http.NewServeMux()
	mux.HandleFunc("/api/puzzles", getPuzzles)

	log.Println("Puzzles API Server is running on port 8081...")
	http.ListenAndServe(":8081", withCORS(mux))
}
