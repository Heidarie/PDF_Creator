namespace PdfKit.Builders;

public sealed record Poster(
    string? Title,
    string? Subtitle,
    string? Location,
    string? Date,
    string? GradientFrom,
    string? GradientTo
);

public static class PosterBuilder
{
    public static Task<string> BuildAsync(Poster poster)
    {
        var model = new
        {
            Title = string.IsNullOrWhiteSpace(poster.Title) ? "Event" : poster.Title,
            poster.Subtitle,
            poster.Location,
            poster.Date,
            GradientFrom = string.IsNullOrWhiteSpace(poster.GradientFrom) ? "#0ea5e9" : poster.GradientFrom,
            GradientTo = string.IsNullOrWhiteSpace(poster.GradientTo) ? "#22c55e" : poster.GradientTo
        };

        const string template = @"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body,html{margin:0;padding:0}
.poster{width:100%;height:100%;min-height:100vh;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,@Model.GradientFrom,@Model.GradientTo);color:white;font-family:Arial,Helvetica,sans-serif}
.card{background:rgba(0,0,0,0.35);padding:40px;border-radius:16px;box-shadow:0 10px 30px rgba(0,0,0,0.35)}
h1{font-size:48px;margin:0}
p{font-size:16px;margin:8px 0 0}
</style>
</head>
<body>
<div class=""poster"">
  <div class=""card"">
    <h1>@Model.Title</h1>
    @if (!string.IsNullOrWhiteSpace(Model.Subtitle))
    {
        <p>@Model.Subtitle</p>
    }
    @if (!string.IsNullOrWhiteSpace(Model.Date) || !string.IsNullOrWhiteSpace(Model.Location))
    {
        <p>@Model.Date@if (!string.IsNullOrWhiteSpace(Model.Date) && !string.IsNullOrWhiteSpace(Model.Location)) { <text> • </text> }@Model.Location</p>
    }
  </div>
</div>
</body>
</html>";

        return TemplateRenderer.RenderAsync("poster", template, model);
    }
}
