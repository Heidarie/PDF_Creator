package builders

import (
	"fmt"
	"math"
	"strconv"
)

type Invoice struct {
	Title    string
	Number   string
	Date     string
	Customer string
	Currency string
	Items    []InvoiceItem
	Notes    string
}

type InvoiceItem struct {
	Description string
	Quantity    float64
	UnitPrice   float64
	Total       float64
}

type invoiceView struct {
	Title    string
	Number   string
	Date     string
	Customer string
	Items    []invoiceLine
	Subtotal string
	Notes    string
}

type invoiceLine struct {
	Description string
	Qty         string
	Unit        string
	Total       string
}

func InvoiceHTML(inv Invoice) (string, error) {
	view := invoiceView{
		Title:    titleOrDefault(inv.Title, "Invoice"),
		Number:   inv.Number,
		Date:     inv.Date,
		Customer: inv.Customer,
		Notes:    inv.Notes,
	}

	currency := currencyOrDefault(inv.Currency)
	subtotal := 0.0

	for _, item := range inv.Items {
		lineTotal := item.Total
		if lineTotal == 0 && item.Quantity > 0 {
			lineTotal = item.Quantity * item.UnitPrice
		}
		subtotal += lineTotal
		view.Items = append(view.Items, invoiceLine{
			Description: item.Description,
			Qty:         formatQty(item.Quantity),
			Unit:        formatMoney(currency, item.UnitPrice),
			Total:       formatMoney(currency, lineTotal),
		})
	}

	view.Subtotal = formatMoney(currency, subtotal)

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
body{font-family:Arial,Helvetica,sans-serif;margin:0;color:#111}
.page{padding:24px}
.header{display:flex;justify-content:space-between;align-items:flex-start}
.title{margin:0 0 8px 0}
.meta{color:#555;font-size:12px}
table{width:100%;border-collapse:collapse;margin-top:16px}
th,td{border-bottom:1px solid #ddd;padding:8px;text-align:left}
th{background:#f6f6f6}
.total{margin-top:16px;text-align:right;font-size:16px;font-weight:bold}
.notes{margin-top:12px;font-size:12px;color:#666}
</style>
</head>
<body>
<div class="page">
  <div class="header">
    <div>
      <h1 class="title">{{ .Title }}</h1>
      {{ if .Customer }}<div class="meta">Customer: {{ .Customer }}</div>{{ end }}
    </div>
    <div class="meta">
      {{ if .Number }}<div><strong>{{ .Number }}</strong></div>{{ end }}
      {{ if .Date }}<div>Date: {{ .Date }}</div>{{ end }}
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th>Item</th>
        <th>Qty</th>
        <th>Unit</th>
        <th>Total</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Items }}
      <tr>
        <td>{{ .Description }}</td>
        <td>{{ .Qty }}</td>
        <td>{{ .Unit }}</td>
        <td>{{ .Total }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>

  <div class="total">Subtotal: {{ .Subtotal }}</div>
  {{ if .Notes }}<div class="notes">{{ .Notes }}</div>{{ end }}
</div>
</body>
</html>`

	return executeTemplate("invoice", tpl, view)
}

func currencyOrDefault(currency string) string {
	if currency == "" {
		return "$"
	}
	return currency
}

func titleOrDefault(title, fallback string) string {
	if title == "" {
		return fallback
	}
	return title
}

func formatQty(qty float64) string {
	if qty == 0 {
		return "0"
	}
	if math.Mod(qty, 1) == 0 {
		return strconv.Itoa(int(qty))
	}
	return fmt.Sprintf("%.2f", qty)
}

func formatMoney(currency string, value float64) string {
	return fmt.Sprintf("%s%.2f", currency, value)
}
