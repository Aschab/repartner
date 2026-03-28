# Pack Calculator

A full-stack application that calculates the optimal pack combination for shipping orders. Given a requested quantity and available pack sizes, it determines which packs to ship while minimizing waste and packaging.

## Overview

This application solves the pack shipping optimization problem with the following priority rules:

1. **Only whole packs can be shipped** - No partial packs allowed
2. **Minimize total items shipped** - Ship as close to the order quantity as possible (always >= order)
3. **Minimize pack count** - When multiple solutions ship the same total, prefer fewer packs

## Architecture

```
pack-calculator/
├── backend/          # Go HTTP API
│   ├── cmd/api/      # Application entrypoint
│   ├── internal/
│   │   ├── config/   # Configuration loading
│   │   ├── domain/   # Domain types
│   │   ├── service/  # Business logic (DP algorithm)
│   │   ├── http/     # HTTP handlers
│   │   └── app/      # Application wiring
│   └── configs/      # Pack size configuration
├── frontend/         # React + TypeScript UI
│   └── src/
│       ├── components/
│       └── api.ts
└── docker-compose.yml
```

## Quick Start with Docker

Everything is dockerized. Just run:

```bash
# Start the application (installs all dependencies automatically)
docker compose up --build

# Or run in detached mode
docker compose up -d --build
```

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

### Custom Ports

If default ports are already in use, you can configure custom ports:

```bash
# Change backend port only
BACKEND_PORT=8081 docker compose up --build -d

# Change frontend port only
FRONTEND_PORT=3000 docker compose up --build -d

# Change both ports
BACKEND_PORT=8081 FRONTEND_PORT=3000 docker compose up --build -d
```

| Variable        | Default | Description                   |
| --------------- | ------- | ----------------------------- |
| `BACKEND_PORT`  | 8080    | Host port for the backend API |
| `FRONTEND_PORT` | 5173    | Host port for the frontend UI |

### Run Tests via Docker

```bash
# Run all tests (backend + frontend)
docker compose --profile test up --build

# Run backend tests only
docker compose build backend-test && docker compose run --rm backend-test

# Run frontend tests only
docker compose build frontend-test && docker compose run --rm frontend-test
```

### Stop Services

```bash
docker compose down
```

## Local Development (Without Docker)

**Backend:**

```bash
cd backend
go mod download        # Install dependencies
go test ./... -v       # Run tests
go run ./cmd/api       # Start server on :8080
```

**Frontend:**

```bash
cd frontend
npm install            # Install dependencies
npm test               # Run tests
npm run dev            # Start dev server on :5173
```

## API Reference

### Health Check

```
GET /health
Response: { "status": "ok" }
```

### Get Pack Sizes

```
GET /api/v1/packs
Response: { "pack_sizes": [250, 500, 1000, 2000, 5000] }
```

### Calculate Packs

```
POST /api/v1/calculate
Content-Type: application/json

Request:  { "order_quantity": 12001 }
Response: {
  "order_quantity": 12001,
  "total_shipped": 12250,
  "total_packs": 4,
  "packs": [
    { "pack_size": 5000, "count": 2 },
    { "pack_size": 2000, "count": 1 },
    { "pack_size": 250, "count": 1 }
  ]
}
```

### Error Response

```json
{
  "error": "order_quantity must be greater than 0"
}
```

## Examples

| Order Qty | Total Shipped | Packs                                                    |
| --------- | ------------- | -------------------------------------------------------- |
| 1         | 250           | 1 × 250                                                  |
| 250       | 250           | 1 × 250                                                  |
| 251       | 500           | 1 × 500                                                  |
| 501       | 750           | 1 × 500, 1 × 250                                         |
| 12001     | 12250         | 2 × 5000, 1 × 2000, 1 × 250                              |
| 500000    | 500000        | 9429 × 53, 7 × 31, 2 × 23 (with pack sizes [23, 31, 53]) |

## Configuration

### Pack Sizes

Pack sizes are configured in `backend/configs/packs.json`:

```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

To change pack sizes:

1. Edit the JSON file
2. Restart the backend service

No code changes required.

### Environment Variables

#### Docker Compose

| Variable        | Default | Description               |
| --------------- | ------- | ------------------------- |
| `BACKEND_PORT`  | 8080    | Host port for backend API |
| `FRONTEND_PORT` | 5173    | Host port for frontend UI |

#### Backend Container

| Variable            | Default            | Description               |
| ------------------- | ------------------ | ------------------------- |
| `PORT`              | 8080               | Internal server port      |
| `PACKS_CONFIG_PATH` | configs/packs.json | Path to pack sizes config |

#### Frontend Build

| Variable            | Default               | Description     |
| ------------------- | --------------------- | --------------- |
| `VITE_API_BASE_URL` | http://localhost:8080 | Backend API URL |

## Algorithm

The calculator uses a dynamic programming approach:

1. Search candidate shipped totals from `orderQty` to `orderQty + maxPack - 1`
2. For each candidate total, use DP to find if it can be formed exactly and with minimum packs
3. Select the smallest feasible total with the fewest packs

This guarantees optimal solutions for any pack size configuration, unlike greedy approaches which can fail for certain pack combinations.

## Tech Stack

- **Backend**: Go 1.22, standard library only
- **Frontend**: React 18, Vite, TypeScript
- **Testing**: Go testing, Vitest + React Testing Library
- **Containerization**: Docker, Docker Compose

## Deployment

### Production URLs

- Frontend: [TBD]
- Backend API: [TBD]

### Deploy to Railway

1. **Push to GitHub**

   ```bash
   git add .
   git commit -m "Ready for deployment"
   git push origin main
   ```

2. **Deploy Backend**

   - Go to [railway.app](https://railway.app)
   - New Project → Deploy from GitHub Repo
   - Select your repo, set root directory to `backend`
   - Railway auto-detects the Dockerfile
   - Note the generated URL (e.g., `https://your-backend.railway.app`)

3. **Deploy Frontend**

   - In the same Railway project, click "New Service"
   - Deploy from GitHub, set root directory to `frontend`
   - Add environment variable: `VITE_API_BASE_URL=https://your-backend.railway.app`
   - Deploy

4. **Update URLs in README** with your deployed URLs

### Deploy to Other Platforms

The Docker images can be deployed to any container platform:

```bash
# Build images
docker compose build

# Push to registry
docker tag repartner-backend:latest your-registry/pack-calculator-backend:latest
docker tag repartner-frontend:latest your-registry/pack-calculator-frontend:latest
docker push your-registry/pack-calculator-backend:latest
docker push your-registry/pack-calculator-frontend:latest
```

## Assumptions & Tradeoffs

1. **Algorithm Choice**: Uses DP instead of greedy to guarantee optimal solutions for any pack configuration
2. **Search Bound**: Searches up to `maxPack - 1` beyond order quantity, sufficient because any larger gap can be filled more efficiently
3. **Memory**: DP uses O(target) space per candidate total; acceptable for typical order sizes
4. **CORS**: Allows all origins for development; should be restricted in production
5. **Validation**: Pack sizes validated at config load and in service layer for defense in depth
