package builders

type Party struct {
	Name  string
	Lines []string
	VAT   string
}

type DetailedInvoice struct {
	CompanyName    string
	CompanyTagline string
	LogoDataURL    string
	Number         string
	Date           string
	DueDate        string
	BillTo         Party
	From           Party
	Currency       string
	Items          []InvoiceItem
	TaxRate        float64
	TaxLabel       string
	Paid           float64
	Notes          string
	Footer         string
	SummaryTitle   string
	SummaryHTML    HTML
	Highlights     []string
}

type detailedInvoiceView struct {
	CompanyName    string
	CompanyTagline string
	LogoDataURL    string
	Number         string
	Date           string
	DueDate        string
	BillTo         Party
	From           Party
	Items          []invoiceLine
	Subtotal       string
	TaxLabel       string
	TaxAmount      string
	Total          string
	Paid           string
	Balance        string
	Notes          string
	Footer         string
	SummaryTitle   string
	SummaryHTML    HTML
	Highlights     []string
	ShowTax        bool
	ShowPaid       bool
}

func DetailedInvoiceHTML(inv DetailedInvoice) (string, error) {
	currency := currencyOrDefault(inv.Currency)
	subtotal := 0.0
	lines := make([]invoiceLine, 0, len(inv.Items))
	for _, item := range inv.Items {
		lineTotal := item.Total
		if lineTotal == 0 && item.Quantity > 0 {
			lineTotal = item.Quantity * item.UnitPrice
		}
		subtotal += lineTotal
		lines = append(lines, invoiceLine{
			Description: item.Description,
			Qty:         formatQty(item.Quantity),
			Unit:        formatMoney(currency, item.UnitPrice),
			Total:       formatMoney(currency, lineTotal),
		})
	}

	taxRate := inv.TaxRate
	if taxRate < 0 {
		taxRate = 0
	}
	if inv.TaxLabel == "" {
		inv.TaxLabel = "Tax"
	}

	taxAmount := subtotal * taxRate
	total := subtotal + taxAmount
	paid := inv.Paid
	balance := total - paid

	view := detailedInvoiceView{
		CompanyName:    inv.CompanyName,
		CompanyTagline: inv.CompanyTagline,
		LogoDataURL:    inv.LogoDataURL,
		Number:         inv.Number,
		Date:           inv.Date,
		DueDate:        inv.DueDate,
		BillTo:         inv.BillTo,
		From:           inv.From,
		Items:          lines,
		Subtotal:       formatMoney(currency, subtotal),
		TaxLabel:       inv.TaxLabel,
		TaxAmount:      formatMoney(currency, taxAmount),
		Total:          formatMoney(currency, total),
		Paid:           formatMoney(currency, paid),
		Balance:        formatMoney(currency, balance),
		Notes:          inv.Notes,
		Footer:         inv.Footer,
		SummaryTitle:   titleOrDefault(inv.SummaryTitle, "Work Summary"),
		SummaryHTML:    inv.SummaryHTML,
		Highlights:     inv.Highlights,
		ShowTax:        taxRate > 0,
		ShowPaid:       paid > 0,
	}

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
body{font-family:Arial,Helvetica,sans-serif;color:#111}
.page{padding:28px}
.row{display:flex;justify-content:space-between;align-items:flex-start}
.logo{display:flex;align-items:center;gap:10px}
.logo img{width:64px;height:64px}
.logo h1{margin:0;font-size:22px}
.meta{font-size:12px;color:#555}
.box{border:1px solid #e5e5e5;border-radius:8px;padding:12px;margin-top:16px}
table{width:100%;border-collapse:collapse;margin-top:12px}
th,td{border-bottom:1px solid #e5e5e5;padding:8px;text-align:left;font-size:12px}
th{background:#f7f7f7}
.right{text-align:right}
.totals{margin-top:16px;display:flex;justify-content:flex-end}
.totals table{width:280px}
.badge{display:inline-block;background:#0ea5e9;color:#fff;padding:4px 8px;border-radius:999px;font-size:11px}
.notes{font-size:11px;color:#666;margin-top:12px}
.footer{margin-top:20px;font-size:11px;color:#777}
.page-break{page-break-after:always}
</style>
</head>
<body>
<div class="page">
  <div class="row">
    <div class="logo">
      {{ if .LogoDataURL }}<img alt="Logo" src="{{ safeURL .LogoDataURL }}">{{ end }}
      <div>
        <h1>{{ if .CompanyName }}{{ .CompanyName }}{{ else }}Company{{ end }}</h1>
        {{ if .CompanyTagline }}<div class="meta">{{ .CompanyTagline }}</div>{{ end }}
      </div>
    </div>
    <div class="meta">
      <div class="badge">INVOICE</div>
      {{ if .Number }}<div><strong># {{ .Number }}</strong></div>{{ end }}
      {{ if .Date }}<div>Date: {{ .Date }}</div>{{ end }}
      {{ if .DueDate }}<div>Due: {{ .DueDate }}</div>{{ end }}
    </div>
  </div>

  <div class="box">
    <div class="row">
      <div>
        <strong>Bill To</strong>
        <div>{{ .BillTo.Name }}</div>
        {{ range .BillTo.Lines }}<div>{{ . }}</div>{{ end }}
        {{ if .BillTo.VAT }}<div>VAT: {{ .BillTo.VAT }}</div>{{ end }}
      </div>
      <div>
        <strong>From</strong>
        <div>{{ .From.Name }}</div>
        {{ range .From.Lines }}<div>{{ . }}</div>{{ end }}
        {{ if .From.VAT }}<div>VAT: {{ .From.VAT }}</div>{{ end }}
      </div>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th style="width:36%">Description</th>
        <th>Qty</th>
        <th class="right">Unit</th>
        <th class="right">Amount</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Items }}
      <tr>
        <td>{{ .Description }}</td>
        <td>{{ .Qty }}</td>
        <td class="right">{{ .Unit }}</td>
        <td class="right">{{ .Total }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>

  <div class="totals">
    <table>
      <tbody>
        <tr><td>Subtotal</td><td class="right">{{ .Subtotal }}</td></tr>
        {{ if .ShowTax }}<tr><td>{{ .TaxLabel }}</td><td class="right">{{ .TaxAmount }}</td></tr>{{ end }}
        <tr><td><strong>Total</strong></td><td class="right"><strong>{{ .Total }}</strong></td></tr>
        {{ if .ShowPaid }}<tr><td>Paid</td><td class="right">{{ .Paid }}</td></tr>{{ end }}
        <tr><td><strong>Balance Due</strong></td><td class="right"><strong>{{ .Balance }}</strong></td></tr>
      </tbody>
    </table>
  </div>

  {{ if .Notes }}<div class="notes">{{ .Notes }}</div>{{ end }}
  {{ if .Footer }}<div class="footer">{{ .Footer }}</div>{{ end }}
</div>

<div class="page-break"></div>

<div class="page">
  <h2>{{ .SummaryTitle }}</h2>
  {{ if .SummaryHTML }}{{ safeHTML .SummaryHTML }}{{ end }}
  {{ if .Highlights }}
  <ul>
    {{ range .Highlights }}<li>{{ . }}</li>{{ end }}
  </ul>
  {{ end }}
</div>
</body>
</html>`

	return executeTemplate("detailed-invoice", tpl, view)
}
