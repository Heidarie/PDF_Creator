using PdfKit.Builders;

namespace PdfKit.Builders.Tests;

public sealed class BuildersTests
{
    [Fact]
    public async Task InvoiceBuilder_RendersTotals()
    {
        var html = await InvoiceBuilder.BuildAsync(new Invoice(
            "Invoice",
            "INV-1",
            "2026-02-14",
            "Acme",
            "$",
            new List<InvoiceItem>
            {
                new("Design", 2, 100),
                new("Audit", 1, 50, 75)
            },
            "Thanks"
        ));

        Assert.Contains("Subtotal", html);
        Assert.Contains("INV-1", html);
    }

    [Fact]
    public async Task DetailedInvoiceBuilder_RendersLogo()
    {
        var html = await DetailedInvoiceBuilder.BuildAsync(new DetailedInvoice(
            "NovaWorks",
            "AI Design Systems",
            "data:image/png;base64,abc",
            "INV-2",
            "2026-02-14",
            "2026-03-01",
            new Party("Acme", new List<string> { "Street 1" }, "VAT1"),
            new Party("Nova", new List<string> { "Street 2" }, "VAT2"),
            "$",
            new List<InvoiceItem> { new("Work", 2, 100) },
            0.2m,
            "Tax",
            50,
            "Pay soon",
            "Footer",
            "Summary",
            "<p>Done</p>",
            new List<string> { "A", "B" }
        ));

        Assert.Contains("data:image/png;base64,abc", html);
        Assert.Contains("Balance Due", html);
    }
}
