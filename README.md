# PDF Creator Service

A tiny Go microservice that renders a PDF from HTML sent over a JSON API. The service exposes a single external endpoint: `POST /render`.

## API

**Request**

`POST /render` with `Content-Type: application/json`

OpenAPI spec: `openapi.yaml` (single external endpoint: `/render`).

```json
{
  "html": "<h1>Invoice</h1><p>Total: $42</p>",
  "name": "invoice-42",
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

**Response**

- `200 OK` with `application/pdf` body
- `400/415/405` for invalid input

**Supported sizes**: `A3`, `A4`, `A5`, `Letter`, `Legal` (case-insensitive).

**Margins** are inches. If omitted, a default of `0.4` in is used.

## Run Locally

```bash
go run ./cmd/pdf-service
```

```bash
curl -X POST http://localhost:8080/render \
  -H 'Content-Type: application/json' \
  -d '{"html":"<h1>Hello</h1>","name":"hello","size":"A4"}' \
  --output hello.pdf
```

## Docker

```bash
docker build -t pdf-service .

docker run --rm -p 8080:8080 pdf-service
```

## Docker Compose

```bash
docker compose up --build
```

## Security Notes

- The renderer blocks all external network requests (`http`, `https`, `file`, `ws`, etc.), preventing SSRF-style access from untrusted HTML.
- Because of this, external images/fonts won't load. Use inline styles or `data:` URLs if you need assets.
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
