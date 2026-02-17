package builders

import (
	"strings"
	"testing"
)

func TestRoadmapHTML(t *testing.T) {
	html, err := RoadmapHTML(Roadmap{
		Title:    "Roadmap",
		Quarters: []string{"Q1", "Q2"},
		Lanes: []RoadmapLane{
			{Label: "Foundation", BarText: "Infra", WidthPct: 30},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "Roadmap") {
		t.Fatalf("expected title")
	}
	if !strings.Contains(html, "Foundation") {
		t.Fatalf("expected lane label")
	}
	if !strings.Contains(html, "background:#3b82f6") {
		t.Fatalf("expected default color")
	}
}
