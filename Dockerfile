FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /out/pdfactory ./cmd/pdfactory

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    ca-certificates \
    fonts-dejavu-core \
    fonts-liberation \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -u 10001 app
USER app
WORKDIR /app
COPY --from=builder /out/pdfactory /app/pdfactory
ENV ADDR=:8080
ENV CHROME_PATH=/usr/bin/chromium
EXPOSE 8080
ENTRYPOINT ["/app/pdfactory"]
