package builders

import (
	"strings"
	"testing"
)

func TestInvoiceHTML(t *testing.T) {
	html, err := InvoiceHTML(Invoice{
		Title:    "Invoice",
		Number:   "INV-1",
		Date:     "2026-02-14",
		Customer: "Acme",
		Currency: "$",
		Items: []InvoiceItem{
			{Description: "Design", Quantity: 2, UnitPrice: 100},
			{Description: "Audit", Quantity: 1, UnitPrice: 50, Total: 75},
		},
		Notes: "Thanks",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "INV-1") {
		t.Fatalf("expected invoice number")
	}
	if !strings.Contains(html, "Subtotal: $") {
		t.Fatalf("expected subtotal")
	}
	if !strings.Contains(html, "Design") {
		t.Fatalf("expected item")
	}
}

func TestDetailedInvoiceHTML(t *testing.T) {
	html, err := DetailedInvoiceHTML(DetailedInvoice{
		CompanyName:    "NovaWorks",
		CompanyTagline: "AI Design Systems",
		LogoDataURL:    "data:image/png;base64,abc",
		Number:         "INV-2",
		Date:           "2026-02-14",
		DueDate:        "2026-03-01",
		BillTo: Party{
			Name:  "Acme",
			Lines: []string{"Street 1"},
			VAT:   "VAT1",
		},
		From: Party{
			Name:  "Nova",
			Lines: []string{"Street 2"},
			VAT:   "VAT2",
		},
		Currency: "$",
		Items: []InvoiceItem{
			{Description: "Work", Quantity: 2, UnitPrice: 100},
		},
		TaxRate:      0.2,
		Paid:         50,
		Notes:        "Pay soon",
		Footer:       "Footer",
		SummaryTitle: "Summary",
		SummaryHTML:  SafeHTML("<p>Done</p>"),
		Highlights:   []string{"A", "B"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "NovaWorks") {
		t.Fatalf("expected company name")
	}
	if !strings.Contains(html, "data:image/png;base64,abc") {
		t.Fatalf("expected logo data url, got: %s", html)
	}
	if !strings.Contains(html, "Balance Due") {
		t.Fatalf("expected balance section")
	}
	if !strings.Contains(html, "Summary") {
		t.Fatalf("expected summary title")
	}
}
