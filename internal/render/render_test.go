package render

import (
	"strings"
	"testing"
	"time"
)

func newTestRenderer(t *testing.T) *Renderer {
	r, err := NewRenderer(Options{
		Timeout:            5 * time.Second,
		MaxHTMLChars:       1000,
		MaxConcurrency:     1,
		DefaultSize:        "A4",
		DefaultOrientation: "portrait",
		DefaultMarginIn:    0.4,
	})
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}
	t.Cleanup(r.Close)
	return r
}

func TestValidateAndNormalizeDefaults(t *testing.T) {
	r := newTestRenderer(t)
	job, err := r.ValidateAndNormalize(Request{HTML: "<p>Hello</p>"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.Filename != "document.pdf" {
		t.Fatalf("expected default filename")
	}
	if !strings.Contains(strings.ToLower(job.HTML), "<html") {
		t.Fatalf("expected html wrapper")
	}
	if job.MarginTop != 0.4 || job.MarginRight != 0.4 || job.MarginBottom != 0.4 || job.MarginLeft != 0.4 {
		t.Fatalf("expected default margins")
	}
}

func TestValidateAndNormalizeErrors(t *testing.T) {
	r := newTestRenderer(t)

	_, err := r.ValidateAndNormalize(Request{HTML: ""})
	if err == nil {
		t.Fatalf("expected error for missing html")
	}

	_, err = r.ValidateAndNormalize(Request{HTML: "x", Size: "A9"})
	if err == nil {
		t.Fatalf("expected error for unsupported size")
	}

	_, err = r.ValidateAndNormalize(Request{HTML: "x", Orientation: "sideways"})
	if err == nil {
		t.Fatalf("expected error for bad orientation")
	}

	neg := -1.0
	_, err = r.ValidateAndNormalize(Request{HTML: "x", Margin: &Margin{Top: &neg}})
	if err == nil {
		t.Fatalf("expected error for negative margin")
	}

	big := 4.0
	_, err = r.ValidateAndNormalize(Request{HTML: "x", Margin: &Margin{Top: &big}})
	if err == nil {
		t.Fatalf("expected error for too-large margin")
	}

	left := 5.0
	right := 5.0
	_, err = r.ValidateAndNormalize(Request{HTML: "x", Margin: &Margin{Left: &left, Right: &right}})
	if err == nil {
		t.Fatalf("expected error for margins too large for page")
	}
}

func TestSanitizeFilename(t *testing.T) {
	name := sanitizeFilename(" ../evil\\name ")
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		t.Fatalf("expected sanitized filename")
	}
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		t.Fatalf("expected .pdf suffix")
	}
}

func TestNormalizeHTML(t *testing.T) {
	wrapped := normalizeHTML("<p>Hi</p>")
	if !strings.Contains(strings.ToLower(wrapped), "<html") {
		t.Fatalf("expected wrapper")
	}

	plain := "<html><body>Hi</body></html>"
	if normalizeHTML(plain) != plain {
		t.Fatalf("expected html to remain unchanged")
	}
}

func TestLookupPaperSize(t *testing.T) {
	if _, ok := lookupPaperSize("letter"); !ok {
		t.Fatalf("expected letter to be supported")
	}
}

func TestBlockedURLPatterns(t *testing.T) {
	patterns := blockedURLPatterns()
	want := []string{"http://*/*", "https://*/*", "file://*/*"}
	for _, p := range want {
		found := false
		for _, actual := range patterns {
			if actual == p {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected pattern %s", p)
		}
	}
}
