# rent-tw-oss

台灣租屋地圖，Go HTTP server + Leaflet.js。

## 執行

```bash
go run main.go
```

開瀏覽器 `http://localhost:8080`

## API

```
GET /api/listings
```

回傳台北、新北、台中、台南、高雄的租屋假資料（JSON）。

## 專案結構

```
.
├── main.go          # HTTP server
├── static/
│   └── index.html   # Leaflet 地圖頁面
├── go.mod
└── .gitignore
```
