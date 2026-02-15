package builders

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
)

type HTML string

func SafeHTML(input string) HTML {
	return HTML(input)
}

func DataURLForSVG(svg []byte) string {
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(svg)
}

func DataURLForPNG(png []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func DataURLForJPEG(jpg []byte) string {
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(jpg)
}

func executeTemplate(name, tpl string, data any) (string, error) {
	funcs := template.FuncMap{
		"safeHTML": func(h HTML) template.HTML {
			return template.HTML(h)
		},
		"safeURL": func(url string) template.URL {
			return template.URL(url)
		},
		"pct": func(v int) string {
			if v < 0 {
				v = 0
			}
			if v > 100 {
				v = 100
			}
			return fmt.Sprintf("%d%%", v)
		},
	}

	t, err := template.New(name).Funcs(funcs).Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
