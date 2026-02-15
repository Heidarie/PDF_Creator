using System.Text;

namespace PdfKit.Builders;

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

public static class ReportBuilder
{
    public static Task<string> BuildAsync(Report report)
    {
        var title = string.IsNullOrWhiteSpace(report.Title) ? "Report" : report.Title;
        var kpis = report.KPIs ?? Array.Empty<KPI>();
        var pages = report.Pages ?? Array.Empty<ReportPage>();

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
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
  <h1>").Append(title).Append(@"</h1>
");
        if (kpis.Count > 0)
        {
            sb.Append("  <div class=\"kpi\">\n");
            foreach (var kpi in kpis)
            {
                sb.Append("    <div class=\"card\">\n");
                sb.Append("      <strong>").Append(kpi.Label).Append("</strong>\n");
                sb.Append("      <div>").Append(kpi.Value).Append("</div>\n");
                sb.Append("      <div class=\"bar\"><span style=\"width:").Append(kpi.Percent).Append("%\"></span></div>\n");
                sb.Append("    </div>\n");
            }
            sb.Append("  </div>\n");
        }
        if (!string.IsNullOrWhiteSpace(report.SummaryHtml))
        {
            sb.Append("  <div class=\"section\">").Append(report.SummaryHtml).Append("</div>\n");
        }
        sb.Append("</div>\n");

        foreach (var page in pages)
        {
            sb.Append("\n<div class=\"page\">\n");
            if (!string.IsNullOrWhiteSpace(page.Title))
            {
                sb.Append("  <h1>").Append(page.Title).Append("</h1>\n");
            }
            if (!string.IsNullOrWhiteSpace(page.BodyHtml))
            {
                sb.Append("  <div class=\"section\">").Append(page.BodyHtml).Append("</div>\n");
            }
            sb.Append("</div>\n");
        }

        if (!string.IsNullOrWhiteSpace(report.AppendixTitle) || !string.IsNullOrWhiteSpace(report.AppendixHtml))
        {
            sb.Append("\n<div>\n");
            if (!string.IsNullOrWhiteSpace(report.AppendixTitle))
            {
                sb.Append("  <h1>").Append(report.AppendixTitle).Append("</h1>\n");
            }
            if (!string.IsNullOrWhiteSpace(report.AppendixHtml))
            {
                sb.Append("  <div class=\"section\">").Append(report.AppendixHtml).Append("</div>\n");
            }
            sb.Append("</div>\n");
        }
        sb.Append(@"</body>
</html>");

        return Task.FromResult(sb.ToString());
    }
}
