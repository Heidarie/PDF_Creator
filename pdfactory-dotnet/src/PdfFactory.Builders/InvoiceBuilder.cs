using System.Globalization;

namespace PdfFactory.Builders;

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

public sealed record SimpleInvoiceLine(string Description, string Qty, string Unit, string Total);

public sealed record InvoiceModel(
    string Title,
    string? Number,
    string? Date,
    string? Customer,
    string? Notes,
    IReadOnlyList<SimpleInvoiceLine> Lines,
    string Subtotal
);

public static class InvoiceBuilder
{
    public static Task<string> BuildAsync(Invoice invoice)
    {
        var currency = string.IsNullOrWhiteSpace(invoice.CurrencySymbol) ? "$" : invoice.CurrencySymbol;
        var lines = invoice.Items.Select(item =>
        {
            var total = item.Total ?? item.UnitPrice * item.Quantity;
            return new SimpleInvoiceLine(
                item.Description,
                FormatQuantity(item.Quantity),
                FormatMoney(currency, item.UnitPrice),
                FormatMoney(currency, total)
            );
        }).ToList();

        var subtotal = invoice.Items.Sum(item => item.Total ?? item.UnitPrice * item.Quantity);

        var model = new InvoiceModel(
            string.IsNullOrWhiteSpace(invoice.Title) ? "Invoice" : invoice.Title,
            invoice.Number,
            invoice.Date,
            invoice.Customer,
            invoice.Notes,
            lines,
            FormatMoney(currency, subtotal)
        );

        const string template = @"<!doctype html>
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
      <h1 class=""title"">@Model.Title</h1>
      @if (!string.IsNullOrWhiteSpace(Model.Customer))
      {
          <div class=""meta"">Customer: @Model.Customer</div>
      }
    </div>
    <div class=""meta"">
      @if (!string.IsNullOrWhiteSpace(Model.Number))
      {
          <div><strong>@Model.Number</strong></div>
      }
      @if (!string.IsNullOrWhiteSpace(Model.Date))
      {
          <div>Date: @Model.Date</div>
      }
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
      @foreach (var line in Model.Lines)
      {
        <tr>
          <td>@line.Description</td>
          <td>@line.Qty</td>
          <td>@line.Unit</td>
          <td>@line.Total</td>
        </tr>
      }
    </tbody>
  </table>

  <div class=""total"">Subtotal: @Model.Subtotal</div>
  @if (!string.IsNullOrWhiteSpace(Model.Notes))
  {
      <div class=""notes"">@Model.Notes</div>
  }
</div>
</body>
</html>";

        return TemplateRenderer.RenderAsync("invoice", template, model);
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
