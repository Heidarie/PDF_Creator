using System.Text.Json;

namespace PdfKit;

public sealed class PdfClientOptions
{
    public string? BaseAddress { get; set; }
    public string? UserAgent { get; set; }
    public long MaxResponseBytes { get; set; } = 25 * 1024 * 1024;
    public JsonSerializerOptions? JsonOptions { get; set; }

    internal const string HttpClientName = "PdfKit";
}
