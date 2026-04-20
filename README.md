# Spend Categoriser

A privacy-first credit card spend categorization tool. Upload PDF credit card statements, get transactions automatically categorized using Merchant Category Codes (MCC) and keyword matching.

<video src="docs/demo.mp4" controls width="100%"></video>

## Privacy

- **All transaction data stays in your browser** (IndexedDB). Nothing is stored on the server.
- The Go backend is stateless — it parses PDFs in memory and returns structured JSON. No data is written to disk, no logs contain PII.
- Docker containers have no volumes. When you stop the app, the server retains zero data.
- No external API calls. Categorization is entirely rule-based and runs locally.

## Quick Start

**Prerequisites:** Docker and Docker Compose

```bash
docker compose up --build
```

Open **http://localhost:3000** in your browser.

To stop:

```bash
docker compose down
```

## Features

- **PDF Upload** — Drag-and-drop one or more credit card PDF statements
- **Auto-Categorization** — Transactions categorized into 11 categories with subcategories:
  - Entertainment, Travel, Utilities, Shopping, Food & Dining, Grocery, Insurance, Healthcare, Education, Lifestyle, Financial
- **Duplicate Detection** — Same transaction won't be counted twice, even if you upload the same bill again (SHA-256 hash of date + description + amount)
- **Dashboard** — Pie chart by category, bar chart of top subcategories, summary cards
- **Transaction List** — Searchable, filterable, sortable table with pagination
- **Bill History** — Track which bills you've uploaded with period and totals
- **Clear Data** — Delete all browser data with one click

## How It Works

1. You upload a PDF credit card statement
2. The Go backend extracts text using `pdftotext` (poppler-utils), identifies transaction lines using regex patterns, and categorizes each transaction using a built-in database of merchant keywords
3. The backend returns structured JSON — the PDF data is not stored anywhere on the server
4. The React frontend deduplicates against existing transactions in IndexedDB and stores new ones
5. Dashboard and charts update automatically

## Architecture

```
Browser (React + Vite)           Docker Network
┌──────────────────────┐        ┌───────────────────────────────────┐
│  Upload UI           │        │  nginx (port 3000)                │
│  Dashboard + Charts  │───────▶│    ├── serves React SPA           │
│  IndexedDB storage   │        │    └── proxies /api/* → backend   │
└──────────────────────┘        │                                   │
                                │  GoFr backend (port 8080)         │
                                │    ├── PDF text extraction         │
                                │    │     (pdftotext / poppler)     │
                                │    ├── Transaction regex parsing   │
                                │    └── MCC keyword categorization  │
                                └───────────────────────────────────┘
```

### Request Flow

```
POST /api/parse (multipart PDF)
  │
  ▼
nginx (:3000)  ──proxy──▶  GoFr (:8080)
                               │
                    ┌──────────┼──────────┐
                    ▼          ▼          ▼
                 parser    extractor  categorizer
               (pdftotext)  (regex)   (keywords)
                    │          │          │
                    └──────────┼──────────┘
                               ▼
                        JSON response
                    (transactions array)
                               │
  ◀────────────────────────────┘
  │
  ▼
Browser: deduplicate → store in IndexedDB → render dashboard
```

## Supported PDF Formats

The parser handles common credit card statement layouts with these date formats:
- `DD/MM/YYYY` (e.g., 01/03/2024)
- `DD-MM-YYYY` (e.g., 01-03-2024)
- `DD Mon YYYY` (e.g., 01 Mar 2024)
- `DD/MM/YY` (e.g., 01/03/24)
- `DD Mon YY` (e.g., 01 Mar 24)

Transaction lines must follow the pattern: **date** → **description** → **amount**

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Health check — returns `{"status": "ok"}` |
| `POST` | `/api/parse` | Upload PDFs (multipart form, field: `files`) |
| `GET` | `/api/categories` | List all categories and subcategories |

### Example: Parse a PDF

```bash
curl -X POST http://localhost:3000/api/parse \
  -F "files=@your_statement.pdf"
```

### Example: List Categories

```bash
curl http://localhost:3000/api/categories
```

## Tech Stack

- **Backend:** Go 1.24, [GoFr](https://gofr.dev/) framework, poppler-utils (`pdftotext`)
- **Frontend:** React 18, TypeScript, Vite 5, Tailwind CSS, Recharts, Dexie.js (IndexedDB)
- **Infrastructure:** Docker multi-stage builds, nginx reverse proxy, Docker Compose

## Project Structure

```
spend-categoriser/
├── docker-compose.yml                  # Service orchestration
├── README.md
│
├── backend/
│   ├── Dockerfile                      # Multi-stage: Go builder → Alpine runtime
│   ├── go.mod                          # Single dependency: gofr.dev
│   ├── configs/
│   │   └── .env                        # GoFr config (port, CORS, logging)
│   ├── cmd/server/
│   │   └── main.go                     # Entry point — route registration
│   └── internal/
│       ├── handler/
│       │   ├── handler.go              # HTTP handlers (Health, ParsePDF, GetCategories)
│       │   ├── errors.go               # Custom error types with HTTP status codes
│       │   └── handler_test.go
│       ├── parser/
│       │   ├── pdf.go                  # PDF → text via pdftotext CLI
│       │   └── pdf_test.go
│       ├── extractor/
│       │   ├── extractor.go            # Regex-based transaction line parsing
│       │   └── extractor_test.go
│       ├── categorizer/
│       │   ├── categorizer.go          # Keyword-based spend categorization
│       │   ├── mcc.go                  # MCC codes + keyword rules database
│       │   └── categorizer_test.go
│       └── models/
│           ├── models.go               # Shared data structures (JSON-serializable)
│           └── models_test.go
│
└── frontend/
    ├── Dockerfile                      # Multi-stage: Node builder → nginx
    ├── nginx.conf                      # Reverse proxy config
    ├── package.json
    ├── vite.config.ts
    ├── tailwind.config.js
    └── src/
        ├── App.tsx                     # Main layout with tab navigation
        ├── main.tsx                    # React entry point
        ├── types/index.ts              # TypeScript interfaces
        ├── db/index.ts                 # IndexedDB schema (Dexie)
        ├── hooks/
        │   └── useTransactions.ts      # Data query hook with aggregations
        ├── utils/
        │   └── dedup.ts                # SHA-256 deduplication logic
        └── components/
            ├── Upload.tsx              # Drag-and-drop PDF upload
            ├── Dashboard.tsx           # Pie/bar charts, summary cards
            ├── TransactionList.tsx     # Filterable, sortable table
            └── BillHistory.tsx         # Uploaded bills tracker
```

## Testing

The backend includes 32 unit tests across all packages:

```bash
# Run tests inside a Go container (pdftotext tests will be skipped)
docker run --rm -v "$(pwd)/backend:/app" -w /app golang:1.24-alpine \
  sh -c "go mod tidy && go test ./... -v"
```

- **categorizer** (6 tests) — Keyword matching, case insensitivity, category validation
- **extractor** (18 tests) — All date formats, credit entries, skip patterns, hash determinism
- **handler** (5 tests) — Health/GetCategories responses, HTTP error status codes
- **models** (3 tests) — JSON field name contracts (guards frontend compatibility)
- **parser** (2 tests) — Invalid/empty input handling (requires `pdftotext`)

## Backend Configuration

GoFr loads configuration from `backend/configs/.env`:

- `APP_NAME=spend-categoriser` — Application name for logging
- `HTTP_PORT=8080` — HTTP server port
- `METRICS_PORT=0` — Prometheus metrics port (disabled)
- `LOG_LEVEL=INFO` — Log verbosity
- `ACCESS_CONTROL_ALLOW_ORIGIN=http://localhost:*` — CORS allowed origins

## Docker Details

The backend runs in a **read-only** container with a tmpfs mount at `/tmp` for security — no persistent writes are possible. GoFr provides built-in structured logging, recovery middleware, and CORS handling.

nginx is configured with:

- 50 MB max upload size for large PDF statements
- 120-second proxy timeouts for processing
- 1-year cache headers for static assets
- SPA fallback routing for React Router
