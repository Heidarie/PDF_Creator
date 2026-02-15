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
        var lines = invoice.Items.Select(item =>
        {
            var total = item.Total ?? item.UnitPrice * item.Quantity;
            return new
            {
                item.Description,
                Qty = InvoiceBuilderHelper.FormatQuantity(item.Quantity),
                Unit = InvoiceBuilderHelper.FormatMoney(currency, item.UnitPrice),
                Total = InvoiceBuilderHelper.FormatMoney(currency, total)
            };
        }).ToList();

        var subtotal = invoice.Items.Sum(item => item.Total ?? item.UnitPrice * item.Quantity);
        var taxRate = invoice.TaxRate < 0 ? 0 : invoice.TaxRate;
        var taxAmount = subtotal * taxRate;
        var totalAmount = subtotal + taxAmount;
        var balance = totalAmount - invoice.Paid;

        var model = new
        {
            CompanyName = invoice.CompanyName,
            CompanyTagline = invoice.CompanyTagline,
            LogoDataUrl = invoice.LogoDataUrl,
            invoice.Number,
            invoice.Date,
            invoice.DueDate,
            BillTo = invoice.BillTo,
            From = invoice.From,
            Lines = lines,
            Subtotal = InvoiceBuilderHelper.FormatMoney(currency, subtotal),
            TaxLabel = string.IsNullOrWhiteSpace(invoice.TaxLabel) ? "Tax" : invoice.TaxLabel,
            TaxAmount = InvoiceBuilderHelper.FormatMoney(currency, taxAmount),
            Total = InvoiceBuilderHelper.FormatMoney(currency, totalAmount),
            Paid = InvoiceBuilderHelper.FormatMoney(currency, invoice.Paid),
            Balance = InvoiceBuilderHelper.FormatMoney(currency, balance),
            Notes = invoice.Notes,
            Footer = invoice.Footer,
            SummaryTitle = string.IsNullOrWhiteSpace(invoice.SummaryTitle) ? "Work Summary" : invoice.SummaryTitle,
            SummaryHtml = invoice.SummaryHtml ?? string.Empty,
            Highlights = invoice.Highlights ?? Array.Empty<string>(),
            ShowTax = taxRate > 0,
            ShowPaid = invoice.Paid > 0
        };

        const string template = @"<!doctype html>
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
      @if (!string.IsNullOrWhiteSpace(Model.LogoDataUrl))
      {
          <img alt=""Logo"" src=""@Model.LogoDataUrl"" />
      }
      <div>
        <h1>@(string.IsNullOrWhiteSpace(Model.CompanyName) ? ""Company"" : Model.CompanyName)</h1>
        @if (!string.IsNullOrWhiteSpace(Model.CompanyTagline))
        {
            <div class=""meta"">@Model.CompanyTagline</div>
        }
      </div>
    </div>
    <div class=""meta"">
      <div class=""badge"">INVOICE</div>
      @if (!string.IsNullOrWhiteSpace(Model.Number))
      {
          <div><strong># @Model.Number</strong></div>
      }
      @if (!string.IsNullOrWhiteSpace(Model.Date))
      {
          <div>Date: @Model.Date</div>
      }
      @if (!string.IsNullOrWhiteSpace(Model.DueDate))
      {
          <div>Due: @Model.DueDate</div>
      }
    </div>
  </div>

  <div class=""box"">
    <div class=""row"">
      <div>
        <strong>Bill To</strong>
        <div>@Model.BillTo.Name</div>
        @foreach (var line in Model.BillTo.Lines)
        {
            <div>@line</div>
        }
        @if (!string.IsNullOrWhiteSpace(Model.BillTo.Vat))
        {
            <div>VAT: @Model.BillTo.Vat</div>
        }
      </div>
      <div>
        <strong>From</strong>
        <div>@Model.From.Name</div>
        @foreach (var line in Model.From.Lines)
        {
            <div>@line</div>
        }
        @if (!string.IsNullOrWhiteSpace(Model.From.Vat))
        {
            <div>VAT: @Model.From.Vat</div>
        }
      </div>
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
      @foreach (var line in Model.Lines)
      {
        <tr>
          <td>@line.Description</td>
          <td>@line.Qty</td>
          <td class=""right"">@line.Unit</td>
          <td class=""right"">@line.Total</td>
        </tr>
      }
    </tbody>
  </table>

  <div class=""totals"">
    <table>
      <tbody>
        <tr><td>Subtotal</td><td class=""right"">@Model.Subtotal</td></tr>
        @if (Model.ShowTax)
        {
            <tr><td>@Model.TaxLabel</td><td class=""right"">@Model.TaxAmount</td></tr>
        }
        <tr><td><strong>Total</strong></td><td class=""right""><strong>@Model.Total</strong></td></tr>
        @if (Model.ShowPaid)
        {
            <tr><td>Paid</td><td class=""right"">@Model.Paid</td></tr>
        }
        <tr><td><strong>Balance Due</strong></td><td class=""right""><strong>@Model.Balance</strong></td></tr>
      </tbody>
    </table>
  </div>

  @if (!string.IsNullOrWhiteSpace(Model.Notes))
  {
      <div class=""notes"">@Model.Notes</div>
  }
  @if (!string.IsNullOrWhiteSpace(Model.Footer))
  {
      <div class=""footer"">@Model.Footer</div>
  }
</div>

<div class=""page-break""></div>

<div class=""page"">
  <h2>@Model.SummaryTitle</h2>
  @Raw(Model.SummaryHtml)
  @if (Model.Highlights.Count > 0)
  {
    <ul>
      @foreach (var item in Model.Highlights)
      {
        <li>@item</li>
      }
    </ul>
  }
</div>
</body>
</html>";

        return TemplateRenderer.RenderAsync("detailed-invoice", template, model);
    }
}

internal static class InvoiceBuilderHelper
{
    public static string FormatMoney(string currency, decimal value) =>
        string.Create(System.Globalization.CultureInfo.InvariantCulture, $"{currency}{value:0.00}");

    public static string FormatQuantity(decimal value)
    {
        if (value == decimal.Truncate(value))
        {
            return ((int)value).ToString(System.Globalization.CultureInfo.InvariantCulture);
        }
        return value.ToString("0.##", System.Globalization.CultureInfo.InvariantCulture);
    }
}
