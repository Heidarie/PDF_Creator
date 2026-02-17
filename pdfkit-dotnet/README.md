# PdfKit (.NET)

A lightweight .NET client and HTML builders for the PDFactory service.

## Projects

- `PdfKit.Client` — calls `/render` and returns PDF bytes
- `PdfKit.Builders` — RazorLight HTML builders (invoice, report, etc.)

## Usage

```csharp
using Microsoft.Extensions.DependencyInjection;
using PdfKit;
using PdfKit.Builders;

var services = new ServiceCollection();
services.AddPdfKitClient(options => options.BaseAddress = "http://localhost:8080/");
var provider = services.BuildServiceProvider();
var client = provider.GetRequiredService<IPdfClient>();

var html = await InvoiceBuilder.BuildAsync(new Invoice(
    "Invoice",
    "INV-2026-014",
    "2026-02-14",
    "Acme Corp",
    "$",
    new List<InvoiceItem> { new("Design", 2, 100) },
    "Thanks for your business!"
));

var req = new RenderRequest { Html = html, Name = "invoice-2026-014", Size = "A4" };
await client.RenderToFileAsync(req, "./outputs/invoice.pdf");
```

## Builders

- `BasicBuilder`
- `InvoiceBuilder`
- `DetailedInvoiceBuilder`
- `ReportBuilder`
- `RoadmapBuilder`
- `PosterBuilder`

## Asset Helpers

Use `DataUrl.FromSvg/FromPng/FromJpeg` to embed images as `data:` URLs. External URLs are blocked by the service.
