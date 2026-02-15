#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cleanup() {
  docker compose -f "${ROOT_DIR}/docker-compose.yml" down -v
}
trap cleanup EXIT

cd "${ROOT_DIR}"

docker compose up -d --build

attempts=0
until [ "$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/render || true)" != "000" ]; do
  attempts=$((attempts + 1))
  if [ "$attempts" -gt 20 ]; then
    echo "Service did not become ready" >&2
    exit 1
  fi
  sleep 1
done

tmp_pdf="${ROOT_DIR}/outputs/ci-hello.pdf"
mkdir -p "${ROOT_DIR}/outputs"

status=$(curl -s -o "${tmp_pdf}" -w "%{http_code}" \
  -X POST http://localhost:8080/render \
  -H 'Content-Type: application/json' \
  -d '{"html":"<h1>Hello</h1><p>PDF from CI</p>","name":"ci-hello","size":"A4"}')

if [ "${status}" != "200" ]; then
  echo "Expected 200, got ${status}" >&2
  echo "Response body:" >&2
  cat "${tmp_pdf}" >&2
  exit 1
fi

if ! head -c 5 "${tmp_pdf}" | grep -q "%PDF-"; then
  echo "Output is not a PDF" >&2
  exit 1
fi

echo "Integration test passed"
