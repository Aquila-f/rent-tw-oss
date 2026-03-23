package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// rawExtracted maps the Python ExtractedRentalFields dataclass.
type rawExtracted struct {
	IsRental          *bool    `json:"is_rental"`
	Title             *string  `json:"title"`
	Price             *int     `json:"price"`
	City              *string  `json:"city"`
	District          *string  `json:"district"`
	Address           *string  `json:"address"`
	AreaPing          *float64 `json:"area_ping"`
	RoomType          *string  `json:"room_type"`
	IsColiving        *bool    `json:"is_coliving"`
	Floor             *int     `json:"floor"`
	TotalFloors       *int     `json:"total_floors"`
	HasElevator       *bool    `json:"has_elevator"`
	HasBalcony        *bool    `json:"has_balcony"`
	HasWashingMachine *bool    `json:"has_washing_machine"`
	PetsAllowed       *bool    `json:"pets_allowed"`
	MinLeaseMonths    *int     `json:"min_lease_months"`
	GenderRestriction *string  `json:"gender_restriction"`
}

// rawRecord maps the Python RentalInfo dataclass (nested JSONL input).
type rawRecord struct {
	PostID    string       `json:"post_id"`
	Content   string       `json:"content"`
	Extracted rawExtracted `json:"extracted"`
	PostURL   *string      `json:"post_url"`
	Timestamp *string      `json:"timestamp"`
	Latitude  *float64     `json:"latitude"`
	Longitude *float64     `json:"longitude"`
	Images    []string     `json:"images"`
}

// RentalPost is the flattened structure served to the frontend.
type RentalPost struct {
	PostID            string   `json:"post_id"`
	PostURL           *string  `json:"post_url"`
	Timestamp         *string  `json:"timestamp"`
	Title             *string  `json:"title"`
	Price             *int     `json:"price"`
	Lat               *float64 `json:"lat"`
	Lng               *float64 `json:"lng"`
	City              *string  `json:"city"`
	District          *string  `json:"district"`
	Address           *string  `json:"address"`
	Content           string   `json:"content"`
	Images            []string `json:"images"`
	AreaPing          *float64 `json:"area_ping"`
	RoomType          *string  `json:"room_type"`
	IsColiving        *bool    `json:"is_coliving"`
	Floor             *int     `json:"floor"`
	TotalFloors       *int     `json:"total_floors"`
	HasElevator       *bool    `json:"has_elevator"`
	HasBalcony        *bool    `json:"has_balcony"`
	HasWashingMachine *bool    `json:"has_washing_machine"`
	PetsAllowed       *bool    `json:"pets_allowed"`
	MinLeaseMonths    *int     `json:"min_lease_months"`
	GenderRestriction *string  `json:"gender_restriction"`
}

func flatten(r rawRecord) RentalPost {
	return RentalPost{
		PostID:            r.PostID,
		PostURL:           r.PostURL,
		Timestamp:         r.Timestamp,
		Title:             r.Extracted.Title,
		Price:             r.Extracted.Price,
		Lat:               r.Latitude,
		Lng:               r.Longitude,
		City:              r.Extracted.City,
		District:          r.Extracted.District,
		Address:           r.Extracted.Address,
		Content:           r.Content,
		Images:            r.Images,
		AreaPing:          r.Extracted.AreaPing,
		RoomType:          r.Extracted.RoomType,
		IsColiving:        r.Extracted.IsColiving,
		Floor:             r.Extracted.Floor,
		TotalFloors:       r.Extracted.TotalFloors,
		HasElevator:       r.Extracted.HasElevator,
		HasBalcony:        r.Extracted.HasBalcony,
		HasWashingMachine: r.Extracted.HasWashingMachine,
		PetsAllowed:       r.Extracted.PetsAllowed,
		MinLeaseMonths:    r.Extracted.MinLeaseMonths,
		GenderRestriction: r.Extracted.GenderRestriction,
	}
}

func loadListings(path string) ([]RentalPost, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := make(map[string]int) // post_id -> index in results
	var results []RentalPost
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var r rawRecord
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			return nil, err
		}
		if r.Extracted.IsRental == nil || !*r.Extracted.IsRental {
			continue
		}
		flat := flatten(r)
		if idx, dup := seen[flat.PostID]; dup {
			results[idx] = flat // keep the latest
		} else {
			seen[flat.PostID] = len(results)
			results = append(results, flat)
		}
	}
	return results, scanner.Err()
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
