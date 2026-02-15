using System.Text;

namespace PdfKit.Builders;

public sealed record RoadmapLane(string Label, string BarText, int WidthPercent, string? Color);
public sealed record Roadmap(string? Title, IReadOnlyList<string> Quarters, IReadOnlyList<RoadmapLane> Lanes);

public static class RoadmapBuilder
{
    public static Task<string> BuildAsync(Roadmap roadmap)
    {
        var title = string.IsNullOrWhiteSpace(roadmap.Title) ? "Product Roadmap" : roadmap.Title;
        var quarters = roadmap.Quarters?.Count > 0 ? roadmap.Quarters : new[] { "Q1", "Q2", "Q3", "Q4" };
        var lanes = roadmap.Lanes ?? Array.Empty<RoadmapLane>();

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body{font-family:Arial,Helvetica,sans-serif;color:#111}
.lane{margin:8px 0}
.bar{display:inline-block;padding:6px 10px;border-radius:4px;color:#fff;font-size:12px}
.grid{display:grid;grid-template-columns:120px 1fr;gap:8px;align-items:center}
.ticks{display:grid;grid-template-columns:repeat(4,1fr);font-size:12px;color:#666;margin-bottom:6px}
</style>
</head>
<body>
<h1>").Append(title).Append(@"</h1>
<div class=""ticks"">
");
        foreach (var q in quarters)
        {
            sb.Append("  <div>").Append(q).Append("</div>\n");
        }
        sb.Append(@"</div>
<div class=""grid"">
");
        foreach (var lane in lanes)
        {
            var color = string.IsNullOrWhiteSpace(lane.Color) ? "#3b82f6" : lane.Color;
            sb.Append("  <div class=\"lane\">").Append(lane.Label).Append("</div>\n");
            sb.Append("  <div class=\"lane\"><span class=\"bar\" style=\"width:")
              .Append(lane.WidthPercent).Append("%;background:").Append(color).Append("\">")
              .Append(lane.BarText).Append("</span></div>\n");
        }
        sb.Append(@"</div>
</body>
</html>");

        return Task.FromResult(sb.ToString());
    }
}
