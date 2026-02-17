package builders

import (
	"strings"
	"testing"
)

func TestBasicHTML(t *testing.T) {
	html, err := BasicHTML(BasicDocument{
		Title: "Hello",
		Body:  SafeHTML("<p>World</p>"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Fatalf("expected title in output")
	}
	if !strings.Contains(html, "<p>World</p>") {
		t.Fatalf("expected body in output")
	}
}
