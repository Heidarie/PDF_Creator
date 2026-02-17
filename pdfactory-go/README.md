# PDFactory Go

A small Go module that provides:

- `pdfactory.Client` for calling the PDFactory service (`POST /render`)
- HTML builders for common documents in `pdfactory-go/builders`

## Install

This module lives as a separate Go module. Update the module path when you publish it.

## Example

```go
package main

import (
	"context"
	"log"

	"pdfactory-go"
	"pdfactory-go/builders"
)

func main() {
	client, err := pdfactory.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	html, err := builders.InvoiceHTML(builders.Invoice{
		Number:   "INV-2026-014",
		Date:     "2026-02-14",
		Customer: "Acme Corp",
		Items: []builders.InvoiceItem{
			{Description: "Design sprint", Quantity: 1, UnitPrice: 2000},
			{Description: "Implementation", Quantity: 5, UnitPrice: 800},
		},
		Notes: "Thanks for your business!",
	})
	if err != nil {
		log.Fatal(err)
	}

	req := pdfactory.RenderRequest{
		HTML: html,
		Name: "invoice-2026-014",
		Size: "A4",
	}

	if err := client.RenderToFile(context.Background(), req, "./outputs/invoice.pdf"); err != nil {
		log.Fatal(err)
	}
}
```

## Builders

- `builders.BasicHTML`
- `builders.InvoiceHTML`
- `builders.DetailedInvoiceHTML`
- `builders.ReportHTML`
- `builders.RoadmapHTML`
- `builders.PosterHTML`

## Asset Helpers

Use `DataURLForSVG/PNG/JPEG` to embed images as `data:` URLs.
