using System.Globalization;
using System.Text;

namespace PdfKit.Builders;

public sealed record InvoiceItem(string Description, decimal Quantity, decimal UnitPrice, decimal? Total = null);

public sealed record Invoice(
    string? Title,
    string? Number,
    string? Date,
    string? Customer,
    string? CurrencySymbol,
    IReadOnlyList<InvoiceItem> Items,
    string? Notes
);

public static class InvoiceBuilder
{
    public static Task<string> BuildAsync(Invoice invoice)
    {
        var currency = string.IsNullOrWhiteSpace(invoice.CurrencySymbol) ? "$" : invoice.CurrencySymbol;
        var title = string.IsNullOrWhiteSpace(invoice.Title) ? "Invoice" : invoice.Title;
        var subtotal = invoice.Items.Sum(item => item.Total ?? item.UnitPrice * item.Quantity);

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
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
<div class=""page"">
  <div class=""header"">
    <div>
      <h1 class=""title"">").Append(title).Append(@"</h1>
");
        if (!string.IsNullOrWhiteSpace(invoice.Customer))
        {
            sb.Append("      <div class=\"meta\">Customer: ").Append(invoice.Customer).Append("</div>\n");
        }
        sb.Append(@"    </div>
    <div class=""meta"">
");
        if (!string.IsNullOrWhiteSpace(invoice.Number))
        {
            sb.Append("      <div><strong>").Append(invoice.Number).Append("</strong></div>\n");
        }
        if (!string.IsNullOrWhiteSpace(invoice.Date))
        {
            sb.Append("      <div>Date: ").Append(invoice.Date).Append("</div>\n");
        }
        sb.Append(@"    </div>
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
");
        foreach (var item in invoice.Items)
        {
            var total = item.Total ?? item.UnitPrice * item.Quantity;
            sb.Append("      <tr>\n");
            sb.Append("        <td>").Append(item.Description).Append("</td>\n");
            sb.Append("        <td>").Append(FormatQuantity(item.Quantity)).Append("</td>\n");
            sb.Append("        <td>").Append(FormatMoney(currency, item.UnitPrice)).Append("</td>\n");
            sb.Append("        <td>").Append(FormatMoney(currency, total)).Append("</td>\n");
            sb.Append("      </tr>\n");
        }
        sb.Append(@"    </tbody>
  </table>

  <div class=""total"">Subtotal: ").Append(FormatMoney(currency, subtotal)).Append(@"</div>
");
        if (!string.IsNullOrWhiteSpace(invoice.Notes))
        {
            sb.Append("  <div class=\"notes\">").Append(invoice.Notes).Append("</div>\n");
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
