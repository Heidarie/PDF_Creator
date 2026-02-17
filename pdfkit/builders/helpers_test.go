package builders

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestDataURLHelpers(t *testing.T) {
	svg := []byte("<svg></svg>")
	svgURL := DataURLForSVG(svg)
	if !strings.HasPrefix(svgURL, "data:image/svg+xml;base64,") {
		t.Fatalf("unexpected svg prefix")
	}
	if !strings.Contains(svgURL, base64.StdEncoding.EncodeToString(svg)) {
		t.Fatalf("svg base64 not found")
	}

	png := []byte("png")
	pngURL := DataURLForPNG(png)
	if !strings.HasPrefix(pngURL, "data:image/png;base64,") {
		t.Fatalf("unexpected png prefix")
	}

	jpg := []byte("jpg")
	jpgURL := DataURLForJPEG(jpg)
	if !strings.HasPrefix(jpgURL, "data:image/jpeg;base64,") {
		t.Fatalf("unexpected jpg prefix")
	}
}
