namespace PdfFactory.Builders;

public sealed record BasicDocument(string? Title, string BodyHtml);

public sealed record BasicDocumentModel(string? Title, string Body);

public static class BasicBuilder
{
    public static Task<string> BuildAsync(BasicDocument doc)
    {
        var body = string.IsNullOrWhiteSpace(doc.BodyHtml) ? "<p>Empty document</p>" : doc.BodyHtml;
        var model = new BasicDocumentModel(
            doc.Title,
            body
        );

        const string template = @"<!doctype html>
<html>
<head>
<meta charset=""utf-8"">
<style>
body{font-family:Arial,Helvetica,sans-serif;padding:24px;color:#111}
</style>
</head>
<body>
@if (!string.IsNullOrWhiteSpace(Model.Title))
{
    <h1>@Model.Title</h1>
}
@Raw(Model.Body)
</body>
</html>";

        return TemplateRenderer.RenderAsync("basic", template, model);
    }
}
