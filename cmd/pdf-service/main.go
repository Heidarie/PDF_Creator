package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"pdf_creator/internal/render"
)

const (
	defaultAddr           = ":8080"
	defaultMaxBodyBytes   = 2 << 20
	defaultRenderTimeout  = 20 * time.Second
	defaultMaxConcurrency = 2
	defaultMaxHTMLChars   = 1_000_000
)

type config struct {
	addr            string
	maxBodyBytes    int64
	renderTimeout   time.Duration
	maxConcurrency  int
	maxHTMLChars    int
	chromePath      string
	chromeNoSandbox bool
}

type renderer interface {
	ValidateAndNormalize(render.Request) (render.Job, error)
	Render(context.Context, render.Job) ([]byte, error)
}

func main() {
	cfg := loadConfig()

	renderer, err := render.NewRenderer(render.Options{
		ChromePath:         cfg.chromePath,
		NoSandbox:          cfg.chromeNoSandbox,
		Timeout:            cfg.renderTimeout,
		MaxHTMLChars:       cfg.maxHTMLChars,
		MaxConcurrency:     cfg.maxConcurrency,
		DefaultSize:        "A4",
		DefaultOrientation: "portrait",
		DefaultMarginIn:    0.4,
	})
	if err != nil {
		log.Fatalf("failed to initialize renderer: %v", err)
	}
	defer renderer.Close()

	handler := &renderHandler{
		renderer:     renderer,
		maxBodyBytes: cfg.maxBodyBytes,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/render", handler.handle)

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           securityHeaders(mux),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("listening on %s", cfg.addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

type renderHandler struct {
	renderer     renderer
	maxBodyBytes int64
}

func (h *renderHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/render" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		writeJSONError(w, http.StatusUnsupportedMediaType, "content-type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodyBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload render.Request
	if err := decoder.Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, jsonErrorMessage(err))
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeJSONError(w, http.StatusBadRequest, "unexpected additional JSON data")
		return
	}

	job, err := h.renderer.ValidateAndNormalize(payload)
	if err != nil {
		var vErr render.ValidationError
		if errors.As(err, &vErr) {
			writeJSONError(w, http.StatusBadRequest, vErr.Error())
			return
		}
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	pdf, err := h.renderer.Render(r.Context(), job)
	if err != nil {
		log.Printf("render error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to render pdf")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+job.Filename+"\"")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdf)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func jsonErrorMessage(err error) string {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxErr):
		return "invalid JSON syntax"
	case errors.As(err, &typeErr):
		if typeErr.Field != "" {
			return "invalid type for field " + typeErr.Field
		}
		return "invalid JSON type"
	case errors.Is(err, io.ErrUnexpectedEOF):
		return "incomplete JSON payload"
	default:
		if strings.HasPrefix(err.Error(), "http: request body too large") {
			return "payload too large"
		}
		return "invalid JSON payload"
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func loadConfig() config {
	cfg := config{
		addr:            getEnv("ADDR", defaultAddr),
		maxBodyBytes:    getEnvInt64("MAX_BODY_BYTES", defaultMaxBodyBytes),
		renderTimeout:   getEnvDuration("RENDER_TIMEOUT", defaultRenderTimeout),
		maxConcurrency:  getEnvInt("MAX_CONCURRENCY", defaultMaxConcurrency),
		maxHTMLChars:    getEnvInt("MAX_HTML_CHARS", defaultMaxHTMLChars),
		chromePath:      os.Getenv("CHROME_PATH"),
		chromeNoSandbox: getEnvBool("CHROME_NO_SANDBOX", false),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.Atoi(val)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if val := os.Getenv(key); val != "" {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		parsed, err := time.ParseDuration(val)
		if err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		switch strings.ToLower(val) {
		case "1", "true", "yes", "y", "on":
			return true
		case "0", "false", "no", "n", "off":
			return false
		}
	}
	return fallback
}
