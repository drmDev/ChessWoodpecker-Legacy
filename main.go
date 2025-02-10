package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// CORS middleware
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Utility function to parse optional integer query params
func parseOptionalInt(param string) *int {
	if param == "" {
		return nil
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return nil // Return nil if parsing fails
	}
	return &value
}

// Get the database connection
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

// Query puzzles from the database
func queryPuzzles(ctx context.Context, conn *pgx.Conn) ([]map[string]interface{}, error) {
	query := `SELECT puzzle_id, category, url FROM puzzles;`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var puzzles []map[string]interface{}
	for rows.Next() {
		var puzzle_id int
		var category, url string
		if err := rows.Scan(&puzzle_id, &category, &url); err != nil {
			return nil, err
		}
		puzzles = append(puzzles, map[string]interface{}{
			"puzzle_id":   puzzle_id,
			"category": category,
			"url":  url,
		})
	}

	return puzzles, nil
}

// Handler function for fetching puzzles
func getPuzzles(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: %s %s", r.Method, r.URL.String())

	conn, err := getDBConnection()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// Query the puzzles from the database
	puzzles, err := queryPuzzles(context.Background(), conn)
	if err != nil {
		log.Printf("Query execution error: %v", err)
		http.Error(w, "Failed to fetch puzzles", http.StatusInternalServerError)
		return
	}

	// Handle no results
	if len(puzzles) == 0 {
		log.Println("No puzzles found.")
		puzzles = []map[string]interface{}{} // Return an empty list
	}

	// Respond with the fetched puzzles
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(puzzles)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/puzzles", getPuzzles) // Expose the /api/puzzles endpoint
	log.Println("Puzzles API Server is running on port 8081...")
	http.ListenAndServe(":8081", withCORS(mux))
}
