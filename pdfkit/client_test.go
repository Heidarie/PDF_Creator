package pdfkit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientRenderSuccess(t *testing.T) {
	pdfBytes := []byte("%PDF-test")

	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			return nil, errors.New("unexpected method")
		}
		if r.URL.Path != "/render" {
			return nil, errors.New("unexpected path")
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return nil, errors.New("missing content-type")
		}
		var req RenderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, errors.New("invalid json")
		}
		if req.HTML == "" {
			return nil, errors.New("empty html")
		}

		header := make(http.Header)
		header.Set("Content-Type", "application/pdf")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(pdfBytes)),
		}, nil
	})

	client, err := NewClient("http://example.test", WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	data, err := client.Render(context.Background(), RenderRequest{HTML: "<h1>hi</h1>"})
	if err != nil {
		t.Fatalf("unexpected render error: %v", err)
	}
	if string(data) != string(pdfBytes) {
		t.Fatalf("unexpected pdf bytes")
	}
}

func TestClientRenderErrorJSON(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, _ := json.Marshal(ErrorResponse{Error: "bad input"})
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(body)),
		}, nil
	})

	client, err := NewClient("http://example.test", WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	_, err = client.Render(context.Background(), RenderRequest{HTML: "<p>x</p>"})
	if err == nil || !strings.Contains(err.Error(), "bad input") {
		t.Fatalf("expected error with message, got: %v", err)
	}
}

func TestClientRenderErrorText(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("boom")),
		}, nil
	})

	client, err := NewClient("http://example.test", WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	_, err = client.Render(context.Background(), RenderRequest{HTML: "<p>x</p>"})
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected error with message, got: %v", err)
	}
}

func TestClientRenderToFile(t *testing.T) {
	pdfBytes := []byte("%PDF-file")
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		header := make(http.Header)
		header.Set("Content-Type", "application/pdf")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(pdfBytes)),
		}, nil
	})

	client, err := NewClient("http://example.test", WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatalf("unexpected client error: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.pdf")
	if err := client.RenderToFile(context.Background(), RenderRequest{HTML: "<p>x</p>"}, outPath); err != nil {
		t.Fatalf("render to file error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if string(data) != string(pdfBytes) {
		t.Fatalf("unexpected file contents")
	}
}
