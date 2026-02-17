using System.Text.Json.Serialization;

namespace PdfKit;

public sealed class RenderRequest
{
    [JsonPropertyName("html")]
    public string Html { get; set; } = string.Empty;

    [JsonPropertyName("name")]
    public string? Name { get; set; }

    [JsonPropertyName("size")]
    public string? Size { get; set; }

    [JsonPropertyName("orientation")]
    public string? Orientation { get; set; }

    [JsonPropertyName("margin")]
    public Margin? Margin { get; set; }
}

public sealed class Margin
{
    [JsonPropertyName("top")]
    public double? Top { get; set; }

    [JsonPropertyName("right")]
    public double? Right { get; set; }

    [JsonPropertyName("bottom")]
    public double? Bottom { get; set; }

    [JsonPropertyName("left")]
    public double? Left { get; set; }
}

public sealed class ErrorResponse
{
    [JsonPropertyName("error")]
    public string? Error { get; set; }
}
