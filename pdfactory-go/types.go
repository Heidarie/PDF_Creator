package pdfactory

type RenderRequest struct {
	HTML        string  `json:"html"`
	Name        string  `json:"name,omitempty"`
	Size        string  `json:"size,omitempty"`
	Orientation string  `json:"orientation,omitempty"`
	Margin      *Margin `json:"margin,omitempty"`
}

type Margin struct {
	Top    *float64 `json:"top,omitempty"`
	Right  *float64 `json:"right,omitempty"`
	Bottom *float64 `json:"bottom,omitempty"`
	Left   *float64 `json:"left,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var SupportedSizes = map[string]struct{}{
	"A3":     {},
	"A4":     {},
	"A5":     {},
	"LETTER": {},
	"LEGAL":  {},
}

var SupportedOrientations = map[string]struct{}{
	"portrait":  {},
	"landscape": {},
}
