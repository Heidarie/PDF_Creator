package builders

import (
	"strings"
	"testing"
)

func TestReportHTML(t *testing.T) {
	html, err := ReportHTML(Report{
		Title:         "Weekly Report",
		Summary:       SafeHTML("<p>Summary text</p>"),
		KPIs:          []KPI{{Label: "Users", Value: "1200", Percent: 75}},
		Pages:         []ReportPage{{Title: "Details", Body: SafeHTML("<p>Details</p>")}},
		AppendixTitle: "Appendix",
		AppendixBody:  SafeHTML("<p>Appendix</p>"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "Weekly Report") {
		t.Fatalf("expected title")
	}
	if !strings.Contains(html, "Users") {
		t.Fatalf("expected KPI")
	}
	if !strings.Contains(html, "Details") {
		t.Fatalf("expected page title")
	}
	if !strings.Contains(html, "Appendix") {
		t.Fatalf("expected appendix")
	}
}
