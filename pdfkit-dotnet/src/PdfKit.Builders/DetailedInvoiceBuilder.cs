using System.Globalization;
using System.Text;

namespace PdfKit.Builders;

public sealed record Party(string Name, IReadOnlyList<string> Lines, string? Vat);

public sealed record DetailedInvoice(
    string? CompanyName,
    string? CompanyTagline,
    string? LogoDataUrl,
    string? Number,
    string? Date,
    string? DueDate,
    Party BillTo,
    Party From,
    string? CurrencySymbol,
    IReadOnlyList<InvoiceItem> Items,
    decimal TaxRate,
    string? TaxLabel,
    decimal Paid,
    string? Notes,
    string? Footer,
    string? SummaryTitle,
    string? SummaryHtml,
    IReadOnlyList<string> Highlights
);

public static class DetailedInvoiceBuilder
{
    public static Task<string> BuildAsync(DetailedInvoice invoice)
    {
        var currency = string.IsNullOrWhiteSpace(invoice.CurrencySymbol) ? "$" : invoice.CurrencySymbol;
        var companyName = string.IsNullOrWhiteSpace(invoice.CompanyName) ? "Company" : invoice.CompanyName;
        var summaryTitle = string.IsNullOrWhiteSpace(invoice.SummaryTitle) ? "Work Summary" : invoice.SummaryTitle;
        var taxLabel = string.IsNullOrWhiteSpace(invoice.TaxLabel) ? "Tax" : invoice.TaxLabel;

        var subtotal = invoice.Items.Sum(item => item.Total ?? item.UnitPrice * item.Quantity);
        var taxRate = invoice.TaxRate < 0 ? 0 : invoice.TaxRate;
        var taxAmount = subtotal * taxRate;
        var totalAmount = subtotal + taxAmount;
        var balance = totalAmount - invoice.Paid;
        var showTax = taxRate > 0;
        var showPaid = invoice.Paid > 0;
        var highlights = invoice.Highlights ?? Array.Empty<string>();

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
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
<div class=""page"">
  <div class=""row"">
    <div class=""logo"">
");
        if (!string.IsNullOrWhiteSpace(invoice.LogoDataUrl))
        {
            sb.Append("      <img alt=\"Logo\" src=\"").Append(invoice.LogoDataUrl).Append("\" />\n");
        }
        sb.Append("      <div>\n");
        sb.Append("        <h1>").Append(companyName).Append("</h1>\n");
        if (!string.IsNullOrWhiteSpace(invoice.CompanyTagline))
        {
            sb.Append("        <div class=\"meta\">").Append(invoice.CompanyTagline).Append("</div>\n");
        }
        sb.Append(@"      </div>
    </div>
    <div class=""meta"">
      <div class=""badge"">INVOICE</div>
");
        if (!string.IsNullOrWhiteSpace(invoice.Number))
        {
            sb.Append("      <div><strong># ").Append(invoice.Number).Append("</strong></div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.Date))
        {
            sb.Append("      <div>Date: ").Append(invoice.Date).Append("</div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.DueDate))
        {
            sb.Append("      <div>Due: ").Append(invoice.DueDate).Append("</div>\n");
        }
        sb.Append(@"    </div>
  </div>

  <div class=""box"">
    <div class=""row"">
      <div>
        <strong>Bill To</strong>
        <div>").Append(invoice.BillTo.Name).Append(@"</div>
");
        foreach (var line in invoice.BillTo.Lines)
        {
            sb.Append("        <div>").Append(line).Append("</div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.BillTo.Vat))
        {
            sb.Append("        <div>VAT: ").Append(invoice.BillTo.Vat).Append("</div>\n");
        }
        sb.Append(@"      </div>
      <div>
        <strong>From</strong>
        <div>").Append(invoice.From.Name).Append(@"</div>
");
        foreach (var line in invoice.From.Lines)
        {
            sb.Append("        <div>").Append(line).Append("</div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.From.Vat))
        {
            sb.Append("        <div>VAT: ").Append(invoice.From.Vat).Append("</div>\n");
        }
        sb.Append(@"      </div>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th style=""width:36%"">Description</th>
        <th>Qty</th>
        <th class=""right"">Unit</th>
        <th class=""right"">Amount</th>
      </tr>
    </thead>
    <tbody>
");
        foreach (var item in invoice.Items)
        {
            var total = item.Total ?? item.UnitPrice * item.Quantity;
            sb.Append("      <tr>\n");
            sb.Append("        <td>").Append(item.Description).Append("</td>\n");
            sb.Append("        <td>").Append(FormatQuantity(item.Quantity)).Append("</td>\n");
            sb.Append("        <td class=\"right\">").Append(FormatMoney(currency, item.UnitPrice)).Append("</td>\n");
            sb.Append("        <td class=\"right\">").Append(FormatMoney(currency, total)).Append("</td>\n");
            sb.Append("      </tr>\n");
        }
        sb.Append(@"    </tbody>
  </table>

  <div class=""totals"">
    <table>
      <tbody>
        <tr><td>Subtotal</td><td class=""right"">").Append(FormatMoney(currency, subtotal)).Append(@"</td></tr>
");
        if (showTax)
        {
            sb.Append("        <tr><td>").Append(taxLabel).Append("</td><td class=\"right\">").Append(FormatMoney(currency, taxAmount)).Append("</td></tr>\n");
        }
        sb.Append("        <tr><td><strong>Total</strong></td><td class=\"right\"><strong>").Append(FormatMoney(currency, totalAmount)).Append("</strong></td></tr>\n");
        if (showPaid)
        {
            sb.Append("        <tr><td>Paid</td><td class=\"right\">").Append(FormatMoney(currency, invoice.Paid)).Append("</td></tr>\n");
        }
        sb.Append("        <tr><td><strong>Balance Due</strong></td><td class=\"right\"><strong>").Append(FormatMoney(currency, balance)).Append(@"</strong></td></tr>
      </tbody>
    </table>
  </div>

");
        if (!string.IsNullOrWhiteSpace(invoice.Notes))
        {
            sb.Append("  <div class=\"notes\">").Append(invoice.Notes).Append("</div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.Footer))
        {
            sb.Append("  <div class=\"footer\">").Append(invoice.Footer).Append("</div>\n");
        }
        sb.Append(@"</div>

<div class=""page-break""></div>

<div class=""page"">
  <h2>").Append(summaryTitle).Append(@"</h2>
");
        if (!string.IsNullOrWhiteSpace(invoice.SummaryHtml))
        {
            sb.Append("  ").Append(invoice.SummaryHtml).Append('\n');
        }
        if (highlights.Count > 0)
        {
            sb.Append("  <ul>\n");
            foreach (var item in highlights)
            {
                sb.Append("    <li>").Append(item).Append("</li>\n");
            }
            sb.Append("  </ul>\n");
        }
        sb.Append(@"</div>
</body>
</html>");

        return Task.FromResult(sb.ToString());
    }

    private static string FormatMoney(string currency, decimal value) =>
        string.Create(CultureInfo.InvariantCulture, $"{currency}{value:0.00}");

    private static string FormatQuantity(decimal value)
    {
        if (value == decimal.Truncate(value))
        {
            return ((int)value).ToString(CultureInfo.InvariantCulture);
        }
        return value.ToString("0.##", CultureInfo.InvariantCulture);
    }
}
