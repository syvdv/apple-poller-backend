# Apple Poller Backend

A lightweight Go backend service that periodically polls the App Store Connect RSS feed to track and store customer reviews for multiple iOS applications.

## Get started

1. Clone the repository :

```bash
git clone https://github.com/syvdv/apple-poller-backend.git
cd apple-poller-backend
```

2. Install Go on your machine following [this link](https://go.dev/learn/).

3. Edit `config.json` if you want to change some parameters.

4. Run this in your terminal to start the server :

```sh
go run ./cmd/poller/main.go
```

## Features

- **Concurrent Architecture** : uses Goroutines to run the background poller and the API server simultaneously.
- **State Persistence** : derives its progress state directly from the data file (`reviews.jsonl`), ensuring no reviews are missed or duplicated upon restart.
- **Auto-Pagination** : automatically fetches up to 10 pages of history to ensure coverage during high-volume periods.
- **JSONL Storage** : uses JSON Lines format for O(1) append performance and memory-efficient scanning.
- **Configurable** : easily adjust polling intervals, app lists, and API settings via `config.json`.

## Project Structure

```text
apple-poller-backend/
├── cmd/
│   └── poller/
│       └── main.go         # Application entry point
├── internal/
│   ├── api/                # HTTP Handlers routing, and middleware
│   ├── config/             # Configuration loading logic
│   ├── fetcher/            # App Store RSS network logic
│   ├── models/             # Data structures
│   └── storage/            # File system read/write logic
├── config.json             # Service configuration
└── reviews.jsonl           # Local data store (generated the app as "DB")
```

## API Endpoints

`GET` /api/apps

```text
Get the list of App IDs available.
```

`GET` /api/reviews?id=####

```text
Get the list of the last 48h reviews for the app with the given id.
```
