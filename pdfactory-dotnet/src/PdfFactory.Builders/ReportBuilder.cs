namespace PdfFactory.Builders;

public sealed record KPI(string Label, string Value, int Percent);
public sealed record ReportPage(string? Title, string BodyHtml);

public sealed record Report(
    string? Title,
    string? SummaryHtml,
    IReadOnlyList<KPI> KPIs,
    IReadOnlyList<ReportPage> Pages,
    string? AppendixTitle,
    string? AppendixHtml
);

public sealed record ReportModel(
    string Title,
    string Summary,
    IReadOnlyList<KPI> KPIs,
    IReadOnlyList<ReportPage> Pages,
    string? AppendixTitle,
    string AppendixHtml
);

public static class ReportBuilder
{
    public static Task<string> BuildAsync(Report report)
    {
        var model = new ReportModel(
            string.IsNullOrWhiteSpace(report.Title) ? "Report" : report.Title,
            report.SummaryHtml ?? string.Empty,
            report.KPIs ?? Array.Empty<KPI>(),
            report.Pages ?? Array.Empty<ReportPage>(),
            report.AppendixTitle,
            report.AppendixHtml ?? string.Empty
        );

        const string template = @"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body{font-family:Georgia,serif;color:#111}
h1{margin-top:0}
.page{page-break-after:always}
.kpi{display:grid;grid-template-columns:repeat(3,1fr);gap:12px}
.card{border:1px solid #ccc;padding:12px;border-radius:6px}
.bar{height:10px;background:#eee;border-radius:6px;overflow:hidden}
.bar>span{display:block;height:10px;background:#3b82f6}
.section{margin-top:16px}
</style>
</head>
<body>
<div class=""page"">
  <h1>@Model.Title</h1>
  @if (Model.KPIs.Count > 0)
  {
    <div class=""kpi"">
      @foreach (var kpi in Model.KPIs)
      {
        <div class=""card"">
          <strong>@kpi.Label</strong>
          <div>@kpi.Value</div>
          <div class=""bar""><span style=""width:@kpi.Percent%""></span></div>
        </div>
      }
    </div>
  }
  @if (!string.IsNullOrWhiteSpace(Model.Summary))
  {
      <div class=""section"">@Raw(Model.Summary)</div>
  }
</div>

@foreach (var page in Model.Pages)
{
<div class=""page"">
  @if (!string.IsNullOrWhiteSpace(page.Title))
  {
      <h1>@page.Title</h1>
  }
  @if (!string.IsNullOrWhiteSpace(page.BodyHtml))
  {
      <div class=""section"">@Raw(page.BodyHtml)</div>
  }
</div>
}

@if (!string.IsNullOrWhiteSpace(Model.AppendixTitle) || !string.IsNullOrWhiteSpace(Model.AppendixHtml))
{
<div>
  @if (!string.IsNullOrWhiteSpace(Model.AppendixTitle))
  {
      <h1>@Model.AppendixTitle</h1>
  }
  @if (!string.IsNullOrWhiteSpace(Model.AppendixHtml))
  {
      <div class=""section"">@Raw(Model.AppendixHtml)</div>
  }
</div>
}
</body>
</html>";

        return TemplateRenderer.RenderAsync("report", template, model);
    }
}
