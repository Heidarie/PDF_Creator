package builders

type Poster struct {
	Title        string
	Subtitle     string
	Location     string
	Date         string
	GradientFrom string
	GradientTo   string
}

type posterView struct {
	Title        string
	Subtitle     string
	Location     string
	Date         string
	GradientFrom string
	GradientTo   string
}

func PosterHTML(poster Poster) (string, error) {
	view := posterView{
		Title:        titleOrDefault(poster.Title, "Event"),
		Subtitle:     poster.Subtitle,
		Location:     poster.Location,
		Date:         poster.Date,
		GradientFrom: colorOrDefault(poster.GradientFrom, "#0ea5e9"),
		GradientTo:   colorOrDefault(poster.GradientTo, "#22c55e"),
	}

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
body,html{margin:0;padding:0}
.poster{width:100%;height:100%;min-height:100vh;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,{{ .GradientFrom }},{{ .GradientTo }});color:white;font-family:Arial,Helvetica,sans-serif}
.card{background:rgba(0,0,0,0.35);padding:40px;border-radius:16px;box-shadow:0 10px 30px rgba(0,0,0,0.35)}
h1{font-size:48px;margin:0}
p{font-size:16px;margin:8px 0 0}
</style>
</head>
<body>
<div class="poster">
  <div class="card">
    <h1>{{ .Title }}</h1>
    {{ if .Subtitle }}<p>{{ .Subtitle }}</p>{{ end }}
    {{ if or .Date .Location }}<p>{{ .Date }}{{ if and .Date .Location }} • {{ end }}{{ .Location }}</p>{{ end }}
  </div>
</div>
</body>
</html>`

	return executeTemplate("poster", tpl, view)
}

func colorOrDefault(color, fallback string) string {
	if color == "" {
		return fallback
	}
	return color
}
