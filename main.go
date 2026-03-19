package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type RentalPost struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Price    int      `json:"price"` // NTD per month
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	City     string   `json:"city"`
	Content  string   `json:"content"`
	Images   []string `json:"images"`
	PostLink string   `json:"post_link"`
}

func loadListings(path string) ([]RentalPost, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var listings []RentalPost
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var p RentalPost
		if err := json.Unmarshal(scanner.Bytes(), &p); err != nil {
			return nil, err
		}
		listings = append(listings, p)
	}
	return listings, scanner.Err()
}

var listings []RentalPost

func listingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(listings)
}

func main() {
	var err error
	listings, err = loadListings("listings.jsonl")
	if err != nil {
		log.Fatalf("failed to load listings: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/listings", listingsHandler)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	addr := ":8080"
	log.Printf("rent-tw listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
