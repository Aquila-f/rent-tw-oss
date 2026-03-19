package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Listing struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price int     `json:"price"` // NTD per month
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	City  string  `json:"city"`
}

var listings = []Listing{
	{1, "大安區溫馨套房", 18000, 25.0330, 121.5654, "台北"},
	{2, "信義區捷運旁一房", 22000, 25.0408, 121.5677, "台北"},
	{3, "中山區明亮雅房", 14500, 25.0630, 121.5248, "台北"},
	{4, "文山區整層分租", 9800, 24.9980, 121.5700, "台北"},
	{5, "板橋兩房寬敞公寓", 16000, 25.0143, 121.4632, "新北"},
	{6, "淡水河景套房", 11000, 25.1683, 121.4381, "新北"},
	{7, "西區中心地帶整層", 12000, 24.1477, 120.6736, "台中"},
	{8, "中興大學周邊雅房", 8500, 24.1185, 120.6843, "台中"},
	{9, "府城歷史街區公寓", 9000, 22.9999, 120.2269, "台南"},
	{10, "愛河旁景觀套房", 13000, 22.6273, 120.3014, "高雄"},
	{11, "左營高鐵旁一房", 10500, 22.6879, 120.2940, "高雄"},
}

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
	mux := http.NewServeMux()
	mux.HandleFunc("/api/listings", listingsHandler)
	mux.Handle("/", http.FileServer(http.Dir("static")))

	addr := ":8080"
	log.Printf("rent-tw listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
