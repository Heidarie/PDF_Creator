package builders

type BasicDocument struct {
	Title string
	Body  HTML
}

func BasicHTML(doc BasicDocument) (string, error) {
	if doc.Body == "" {
		doc.Body = HTML("<p>Empty document</p>")
	}

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
body{font-family:Arial,Helvetica,sans-serif;padding:24px;color:#111}
</style>
</head>
<body>
{{ if .Title }}<h1>{{ .Title }}</h1>{{ end }}
{{ safeHTML .Body }}
</body>
</html>`

	return executeTemplate("basic", tpl, doc)
}
