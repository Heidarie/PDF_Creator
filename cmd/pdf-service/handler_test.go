package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pdf_creator/internal/render"
)

type fakeRenderer struct {
	validate func(render.Request) (render.Job, error)
	render   func(context.Context, render.Job) ([]byte, error)
}

func (f fakeRenderer) ValidateAndNormalize(req render.Request) (render.Job, error) {
	if f.validate != nil {
		return f.validate(req)
	}
	return render.Job{}, nil
}

func (f fakeRenderer) Render(ctx context.Context, job render.Job) ([]byte, error) {
	if f.render != nil {
		return f.render(ctx, job)
	}
	return nil, nil
}

func TestHandleMethodNotAllowed(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{}, maxBodyBytes: 1024}
	req := httptest.NewRequest(http.MethodGet, "/render", nil)
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Allow") != http.MethodPost {
		t.Fatalf("expected Allow header")
	}
}

func TestHandleUnsupportedMediaType(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{}, maxBodyBytes: 1024}
	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", resp.StatusCode)
	}
}

func TestHandleInvalidJSON(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{}, maxBodyBytes: 1024}
	req := httptest.NewRequest(http.MethodPost, "/render", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var out map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out["error"] == "" {
		t.Fatalf("expected error message")
	}
}

func TestHandleUnknownFields(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{}, maxBodyBytes: 1024}
	body := bytes.NewBufferString(`{"html":"x","extra":1}`)
	req := httptest.NewRequest(http.MethodPost, "/render", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandleValidationError(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{
		validate: func(render.Request) (render.Job, error) {
			return render.Job{}, render.ValidationError{Field: "html", Message: "is required"}
		},
	}, maxBodyBytes: 1024}

	body := bytes.NewBufferString(`{"html":""}`)
	req := httptest.NewRequest(http.MethodPost, "/render", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestHandleRenderError(t *testing.T) {
	h := &renderHandler{renderer: fakeRenderer{
		validate: func(render.Request) (render.Job, error) {
			return render.Job{Filename: "test.pdf"}, nil
		},
		render: func(context.Context, render.Job) ([]byte, error) {
			return nil, errors.New("boom")
		},
	}, maxBodyBytes: 1024}

	body := bytes.NewBufferString(`{"html":"<p>x</p>"}`)
	req := httptest.NewRequest(http.MethodPost, "/render", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

func TestHandleSuccess(t *testing.T) {
	pdf := []byte("%PDF-1.4")
	h := &renderHandler{renderer: fakeRenderer{
		validate: func(render.Request) (render.Job, error) {
			return render.Job{Filename: "ok.pdf"}, nil
		},
		render: func(context.Context, render.Job) ([]byte, error) {
			return pdf, nil
		},
	}, maxBodyBytes: 1024}

	body := bytes.NewBufferString(`{"html":"<p>x</p>"}`)
	req := httptest.NewRequest(http.MethodPost, "/render", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("expected pdf content-type")
	}
	if disp := resp.Header.Get("Content-Disposition"); !strings.Contains(disp, "ok.pdf") {
		t.Fatalf("expected content-disposition filename")
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if !bytes.Equal(bodyBytes, pdf) {
		t.Fatalf("unexpected body")
	}
}

func TestSecurityHeaders(t *testing.T) {
	h := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Header.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("missing X-Frame-Options")
	}
	if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing X-Content-Type-Options")
	}
	if resp.Header.Get("Referrer-Policy") != "no-referrer" {
		t.Fatalf("missing Referrer-Policy")
	}
}
