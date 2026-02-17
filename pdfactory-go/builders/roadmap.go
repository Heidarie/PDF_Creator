package builders

type Roadmap struct {
	Title    string
	Quarters []string
	Lanes    []RoadmapLane
}

type RoadmapLane struct {
	Label    string
	BarText  string
	WidthPct int
	Color    string
}

type roadmapView struct {
	Title    string
	Quarters []string
	Lanes    []RoadmapLane
}

func RoadmapHTML(roadmap Roadmap) (string, error) {
	if len(roadmap.Quarters) == 0 {
		roadmap.Quarters = []string{"Q1", "Q2", "Q3", "Q4"}
	}
	lanes := make([]RoadmapLane, 0, len(roadmap.Lanes))
	for _, lane := range roadmap.Lanes {
		if lane.Color == "" {
			lane.Color = "#3b82f6"
		}
		lanes = append(lanes, lane)
	}
	view := roadmapView{
		Title:    titleOrDefault(roadmap.Title, "Product Roadmap"),
		Quarters: roadmap.Quarters,
		Lanes:    lanes,
	}

	const tpl = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<style>
body{font-family:Arial,Helvetica,sans-serif;color:#111}
.lane{margin:8px 0}
.bar{display:inline-block;padding:6px 10px;border-radius:4px;color:#fff;font-size:12px}
.grid{display:grid;grid-template-columns:120px 1fr;gap:8px;align-items:center}
.ticks{display:grid;grid-template-columns:repeat(4,1fr);font-size:12px;color:#666;margin-bottom:6px}
</style>
</head>
<body>
<h1>{{ .Title }}</h1>
<div class="ticks">
  {{ range .Quarters }}<div>{{ . }}</div>{{ end }}
</div>
<div class="grid">
  {{ range .Lanes }}
  <div class="lane">{{ .Label }}</div>
  <div class="lane"><span class="bar" style="width:{{ pct .WidthPct }};background:{{ .Color }}">{{ .BarText }}</span></div>
  {{ end }}
</div>
</body>
</html>`

	return executeTemplate("roadmap", tpl, view)
}
