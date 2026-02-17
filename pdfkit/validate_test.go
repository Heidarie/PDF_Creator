package pdfkit

import "testing"

func TestValidateRequest(t *testing.T) {
	if err := ValidateRequest(RenderRequest{}); err == nil {
		t.Fatalf("expected error for empty html")
	}

	ok := RenderRequest{HTML: "<h1>hi</h1>"}
	if err := ValidateRequest(ok); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	badSize := RenderRequest{HTML: "x", Size: "A9"}
	if err := ValidateRequest(badSize); err == nil {
		t.Fatalf("expected error for size")
	}

	badOrientation := RenderRequest{HTML: "x", Orientation: "sideways"}
	if err := ValidateRequest(badOrientation); err == nil {
		t.Fatalf("expected error for orientation")
	}

	caseOk := RenderRequest{HTML: "x", Size: "letter", Orientation: "LANDSCAPE"}
	if err := ValidateRequest(caseOk); err != nil {
		t.Fatalf("unexpected error for case-insensitive values: %v", err)
	}
}
