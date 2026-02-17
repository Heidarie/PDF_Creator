# PDFactory

PDFactory is a lightweight microservice that renders PDFs from HTML sent in JSON. It exposes a single external endpoint and is designed to be container‑first and SSRF‑safe.

## Highlights

- Single `POST /render` endpoint
- Headless Chromium rendering (Docker ready)
- SSRF‑safe: blocks external network requests from HTML
- Strict input validation and size limits
- Go client + HTML builders (`pdfactory-go/`)
- .NET client + RazorLight builders (`pdfactory-dotnet/`)

## Quickstart (Docker)

```bash
docker build -t pdfactory .

docker run --rm -p 8080:8080 pdfactory
```

Test it:

```bash
curl -X POST http://localhost:8080/render \
  -H 'Content-Type: application/json' \
  -d '{"html":"<h1>Hello</h1><p>PDF from PDFactory</p>","name":"hello","size":"A4"}' \
  --output hello.pdf
```

## API

`POST /render` with `Content-Type: application/json`

```json
{
  "html": "<h1>Hello</h1><p>PDF from PDFactory</p>",
  "name": "hello",
  "size": "A4",
  "orientation": "portrait",
  "margin": {
    "top": 0.4,
    "right": 0.4,
    "bottom": 0.4,
    "left": 0.4
  }
}
```

Response:

- `200 OK` with `application/pdf` body
- `400/415/405` for invalid input

Supported sizes (case‑insensitive): `A3`, `A4`, `A5`, `Letter`, `Legal`

Margins are inches. Defaults to `0.4` if omitted.

OpenAPI spec: `openapi.yaml` (single external endpoint: `/render`).

## Security Notes

- External requests are blocked (`http`, `https`, `file`, `ws`, etc.) to prevent SSRF.
- Use inline CSS and `data:` URLs for images/fonts.
- Request body size, HTML length, and concurrency are bounded.

## Configuration

Environment variables:

- `ADDR` (default `:8080`)
- `MAX_BODY_BYTES` (default `2097152`)
- `MAX_HTML_CHARS` (default `1000000`)
- `MAX_CONCURRENCY` (default `2`)
- `RENDER_TIMEOUT` (default `20s`)
- `CHROME_PATH` (default `/usr/bin/chromium` inside Docker)
- `CHROME_NO_SANDBOX` (default `false`) — set to `true` only if Chromium fails to start in your container

## Docker Compose

```bash
docker compose up --build
```

Optional Redis profile:

```bash
docker compose --profile infra up --build
```

## Clients

Go client + builders:

- `pdfactory-go/`

.NET client + RazorLight builders:

- `pdfactory-dotnet/`

## Testing

Go tests:

```bash
go test ./...
```

Go client tests:

```bash
cd pdfactory-go
go test ./...
```

Integration tests (Docker):

```bash
bash scripts/integration-test.sh
```

.NET tests:

```bash
cd pdfactory-dotnet
dotnet test
```

## Project Layout

- `cmd/pdfactory/` — HTTP server
- `internal/render/` — Chromium rendering + validation
- `pdfactory-go/` — Go client + HTML builders
- `pdfactory-dotnet/` — .NET client + RazorLight builders
- `scripts/` — CI and integration scripts
