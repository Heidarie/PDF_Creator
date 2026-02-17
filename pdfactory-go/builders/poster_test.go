package builders

import (
	"strings"
	"testing"
)

func TestPosterHTML(t *testing.T) {
	html, err := PosterHTML(Poster{
		Title:    "Design Summit",
		Subtitle: "Future",
		Date:     "2026-03-01",
		Location: "Warsaw",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "Design Summit") {
		t.Fatalf("expected title")
	}
	if !strings.Contains(html, "Warsaw") {
		t.Fatalf("expected location")
	}
	if !strings.Contains(html, "linear-gradient") {
		t.Fatalf("expected gradient")
	}
}
