package builders

type KPI struct {
	Label   string
	Value   string
	Percent int
}

type ReportPage struct {
	Title string
	Body  HTML
}

type Report struct {
	Title         string
	Summary       HTML
	KPIs          []KPI
	Pages         []ReportPage
	AppendixTitle string
	AppendixBody  HTML
}

type reportPageView struct {
	Title     string
	Body      HTML
	PageBreak bool
}

type reportView struct {
	Title         string
	Summary       HTML
	KPIs          []KPI
	Pages         []reportPageView
	AppendixTitle string
	AppendixBody  HTML
}

func ReportHTML(report Report) (string, error) {
	pages := make([]reportPageView, 0, len(report.Pages))
	for i, page := range report.Pages {
		pages = append(pages, reportPageView{
			Title:     page.Title,
			Body:      page.Body,
			PageBreak: i < len(report.Pages)-1 || report.AppendixTitle != "" || report.AppendixBody != "",
		})
	}

	view := reportView{
		Title:         titleOrDefault(report.Title, "Report"),
		Summary:       report.Summary,
		KPIs:          report.KPIs,
		Pages:         pages,
		AppendixTitle: report.AppendixTitle,
		AppendixBody:  report.AppendixBody,
	}

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
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
<div class="page">
  <h1>{{ .Title }}</h1>
  {{ if .KPIs }}
  <div class="kpi">
    {{ range .KPIs }}
    <div class="card">
      <strong>{{ .Label }}</strong>
      <div>{{ .Value }}</div>
      <div class="bar"><span style="width:{{ pct .Percent }}"></span></div>
    </div>
    {{ end }}
  </div>
  {{ end }}
  {{ if .Summary }}<div class="section">{{ safeHTML .Summary }}</div>{{ end }}
</div>

{{ range .Pages }}
<div{{ if .PageBreak }} class="page"{{ end }}>
  {{ if .Title }}<h1>{{ .Title }}</h1>{{ end }}
  {{ if .Body }}<div class="section">{{ safeHTML .Body }}</div>{{ end }}
</div>
{{ end }}

{{ if or .AppendixTitle .AppendixBody }}
<div>
  {{ if .AppendixTitle }}<h1>{{ .AppendixTitle }}</h1>{{ end }}
  {{ if .AppendixBody }}<div class="section">{{ safeHTML .AppendixBody }}</div>{{ end }}
</div>
{{ end }}
</body>
</html>`

	return executeTemplate("report", tpl, view)
}
