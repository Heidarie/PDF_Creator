using System.Text;

namespace PdfKit.Builders;

public sealed record BasicDocument(string? Title, string BodyHtml);

public static class BasicBuilder
{
    public static Task<string> BuildAsync(BasicDocument doc)
    {
        var body = string.IsNullOrWhiteSpace(doc.BodyHtml) ? "<p>Empty document</p>" : doc.BodyHtml;

        var sb = new StringBuilder();
        sb.Append(@"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body{font-family:Arial,Helvetica,sans-serif;padding:24px;color:#111}
</style>
</head>
<body>
");
        if (!string.IsNullOrWhiteSpace(doc.Title))
        {
            sb.Append("<h1>").Append(doc.Title).Append("</h1>\n");
        }
        sb.Append(body);
        sb.Append(@"
</body>
</html>");

        return Task.FromResult(sb.ToString());
    }
}
