namespace PdfKit.Builders;

public sealed record RoadmapLane(string Label, string BarText, int WidthPercent, string? Color);
public sealed record Roadmap(string? Title, IReadOnlyList<string> Quarters, IReadOnlyList<RoadmapLane> Lanes);

internal sealed record RoadmapLaneView(string Label, string BarText, int WidthPercent, string Color);

public static class RoadmapBuilder
{
    public static Task<string> BuildAsync(Roadmap roadmap)
    {
        var quarters = roadmap.Quarters?.Count > 0 ? roadmap.Quarters : new[] { "Q1", "Q2", "Q3", "Q4" };
        var lanes = roadmap.Lanes?.Select(lane => new RoadmapLaneView(
            lane.Label,
            lane.BarText,
            lane.WidthPercent,
            string.IsNullOrWhiteSpace(lane.Color) ? "#3b82f6" : lane.Color
        )).ToList() ?? new List<RoadmapLaneView>();

        var model = new
        {
            Title = string.IsNullOrWhiteSpace(roadmap.Title) ? "Product Roadmap" : roadmap.Title,
            Quarters = quarters,
            Lanes = lanes
        };

        const string template = @"<!doctype html>
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
<h1>@Model.Title</h1>
<div class=""ticks"">
  @foreach (var q in Model.Quarters)
  {
      <div>@q</div>
  }
</div>
<div class=""grid"">
  @foreach (var lane in Model.Lanes)
  {
      <div class=""lane"">@lane.Label</div>
      <div class=""lane""><span class=""bar"" style=""width:@lane.WidthPercent%;background:@lane.Color"">@lane.BarText</span></div>
  }
</div>
</body>
</html>";

        return TemplateRenderer.RenderAsync("roadmap", template, model);
    }
}
