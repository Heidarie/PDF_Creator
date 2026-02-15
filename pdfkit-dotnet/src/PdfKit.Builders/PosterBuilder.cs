using System.Text;

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
        var title = string.IsNullOrWhiteSpace(poster.Title) ? "Event" : poster.Title;
        var gradientFrom = string.IsNullOrWhiteSpace(poster.GradientFrom) ? "#0ea5e9" : poster.GradientFrom;
        var gradientTo = string.IsNullOrWhiteSpace(poster.GradientTo) ? "#22c55e" : poster.GradientTo;

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body,html{margin:0;padding:0}
.poster{width:100%;height:100%;min-height:100vh;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,").Append(gradientFrom).Append(',').Append(gradientTo).Append(@");color:white;font-family:Arial,Helvetica,sans-serif}
.card{background:rgba(0,0,0,0.35);padding:40px;border-radius:16px;box-shadow:0 10px 30px rgba(0,0,0,0.35)}
h1{font-size:48px;margin:0}
p{font-size:16px;margin:8px 0 0}
</style>
</head>
<body>
<div class=""poster"">
  <div class=""card"">
    <h1>").Append(title).Append(@"</h1>
");
        if (!string.IsNullOrWhiteSpace(poster.Subtitle))
        {
            sb.Append("    <p>").Append(poster.Subtitle).Append("</p>\n");
        }
        if (!string.IsNullOrWhiteSpace(poster.Date) || !string.IsNullOrWhiteSpace(poster.Location))
        {
            sb.Append("    <p>");
            if (!string.IsNullOrWhiteSpace(poster.Date))
            {
                sb.Append(poster.Date);
            }
            if (!string.IsNullOrWhiteSpace(poster.Date) && !string.IsNullOrWhiteSpace(poster.Location))
            {
                sb.Append(" • ");
            }
            if (!string.IsNullOrWhiteSpace(poster.Location))
            {
                sb.Append(poster.Location);
            }
            sb.Append("</p>\n");
        }
        sb.Append(@"  </div>
</div>
</body>
</html>");

        return Task.FromResult(sb.ToString());
    }
}
