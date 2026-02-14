package render

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Options struct {
	ChromePath         string
	Timeout            time.Duration
	MaxHTMLChars       int
	MaxConcurrency     int
	DefaultSize        string
	DefaultOrientation string
	DefaultMarginIn    float64
	NoSandbox          bool
}

type Request struct {
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

type Job struct {
	HTML         string
	Filename     string
	PaperWidth   float64
	PaperHeight  float64
	Landscape    bool
	MarginTop    float64
	MarginRight  float64
	MarginBottom float64
	MarginLeft   float64
}

type Renderer struct {
	allocCtx    context.Context
	allocCancel context.CancelFunc
	opts        Options
	sem         chan struct{}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func NewRenderer(opts Options) (*Renderer, error) {
	if opts.MaxConcurrency <= 0 {
		opts.MaxConcurrency = 1
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 20 * time.Second
	}
	if opts.MaxHTMLChars <= 0 {
		opts.MaxHTMLChars = 1_000_000
	}
	if opts.DefaultSize == "" {
		opts.DefaultSize = "A4"
	}
	if opts.DefaultOrientation == "" {
		opts.DefaultOrientation = "portrait"
	}
	if opts.DefaultMarginIn <= 0 {
		opts.DefaultMarginIn = 0.4
	}

	allocOpts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-domain-reliability", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	}
	if opts.NoSandbox {
		allocOpts = append(allocOpts,
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-setuid-sandbox", true),
		)
	}
	if opts.ChromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(opts.ChromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	return &Renderer{
		allocCtx:    allocCtx,
		allocCancel: allocCancel,
		opts:        opts,
		sem:         make(chan struct{}, opts.MaxConcurrency),
	}, nil
}

func (r *Renderer) Close() {
	if r.allocCancel != nil {
		r.allocCancel()
	}
}

func (r *Renderer) ValidateAndNormalize(req Request) (Job, error) {
	html := strings.TrimSpace(req.HTML)
	if html == "" {
		return Job{}, ValidationError{Field: "html", Message: "is required"}
	}
	if len(html) > r.opts.MaxHTMLChars {
		return Job{}, ValidationError{Field: "html", Message: "is too large"}
	}

	name := strings.TrimSpace(req.Name)
	filename := sanitizeFilename(name)

	sizeName := strings.TrimSpace(req.Size)
	if sizeName == "" {
		sizeName = r.opts.DefaultSize
	}

	size, ok := lookupPaperSize(sizeName)
	if !ok {
		return Job{}, ValidationError{Field: "size", Message: fmt.Sprintf("unsupported size %q", sizeName)}
	}

	orientation := strings.TrimSpace(req.Orientation)
	if orientation == "" {
		orientation = r.opts.DefaultOrientation
	}
	orientation = strings.ToLower(orientation)

	landscape := false
	switch orientation {
	case "portrait":
		landscape = false
	case "landscape":
		landscape = true
	default:
		return Job{}, ValidationError{Field: "orientation", Message: "must be portrait or landscape"}
	}

	marginTop := r.opts.DefaultMarginIn
	marginRight := r.opts.DefaultMarginIn
	marginBottom := r.opts.DefaultMarginIn
	marginLeft := r.opts.DefaultMarginIn
	if req.Margin != nil {
		if req.Margin.Top != nil {
			marginTop = *req.Margin.Top
		}
		if req.Margin.Right != nil {
			marginRight = *req.Margin.Right
		}
		if req.Margin.Bottom != nil {
			marginBottom = *req.Margin.Bottom
		}
		if req.Margin.Left != nil {
			marginLeft = *req.Margin.Left
		}
	}

	const maxMarginIn = 3.0
	if marginTop < 0 || marginRight < 0 || marginBottom < 0 || marginLeft < 0 {
		return Job{}, ValidationError{Field: "margin", Message: "must be non-negative"}
	}
	if marginTop > maxMarginIn || marginRight > maxMarginIn || marginBottom > maxMarginIn || marginLeft > maxMarginIn {
		return Job{}, ValidationError{Field: "margin", Message: "is too large"}
	}

	width := size.Width
	height := size.Height
	if landscape {
		width, height = height, width
	}
	if marginLeft+marginRight >= width-0.1 {
		return Job{}, ValidationError{Field: "margin", Message: "left+right margins are too large for the page"}
	}
	if marginTop+marginBottom >= height-0.1 {
		return Job{}, ValidationError{Field: "margin", Message: "top+bottom margins are too large for the page"}
	}

	normalizedHTML := normalizeHTML(html)

	return Job{
		HTML:         normalizedHTML,
		Filename:     filename,
		PaperWidth:   size.Width,
		PaperHeight:  size.Height,
		Landscape:    landscape,
		MarginTop:    marginTop,
		MarginRight:  marginRight,
		MarginBottom: marginBottom,
		MarginLeft:   marginLeft,
	}, nil
}

func (r *Renderer) Render(reqCtx context.Context, job Job) ([]byte, error) {
	if err := r.acquire(reqCtx); err != nil {
		return nil, err
	}
	defer r.release()

	taskCtx, taskCancel := chromedp.NewContext(r.allocCtx)
	defer taskCancel()

	renderCtx, renderCancel := context.WithTimeout(taskCtx, r.opts.Timeout)
	defer renderCancel()

	go func() {
		select {
		case <-reqCtx.Done():
			renderCancel()
		case <-renderCtx.Done():
		}
	}()

	dataURL := "data:text/html;base64," + base64.StdEncoding.EncodeToString([]byte(job.HTML))

	var pdf []byte
	actions := []chromedp.Action{
		network.Enable(),
		network.SetBlockedURLS(blockedURLPatterns()),
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(job.PaperWidth).
				WithPaperHeight(job.PaperHeight).
				WithLandscape(job.Landscape).
				WithMarginTop(job.MarginTop).
				WithMarginBottom(job.MarginBottom).
				WithMarginLeft(job.MarginLeft).
				WithMarginRight(job.MarginRight).
				Do(ctx)
			if err != nil {
				return err
			}
			pdf = buf
			return nil
		}),
	}

	if err := chromedp.Run(renderCtx, actions...); err != nil {
		return nil, err
	}

	if len(pdf) == 0 {
		return nil, errors.New("empty pdf output")
	}

	return pdf, nil
}

func (r *Renderer) acquire(ctx context.Context) error {
	select {
	case r.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Renderer) release() {
	select {
	case <-r.sem:
	default:
	}
}

type paperSize struct {
	Width  float64
	Height float64
}

var paperSizes = map[string]paperSize{
	"A3":     {Width: 11.69, Height: 16.54},
	"A4":     {Width: 8.27, Height: 11.69},
	"A5":     {Width: 5.83, Height: 8.27},
	"LETTER": {Width: 8.5, Height: 11},
	"LEGAL":  {Width: 8.5, Height: 14},
}

func lookupPaperSize(name string) (paperSize, bool) {
	key := strings.ToUpper(strings.TrimSpace(name))
	size, ok := paperSizes[key]
	return size, ok
}

func normalizeHTML(html string) string {
	lower := strings.ToLower(html)
	if strings.Contains(lower, "<html") {
		return html
	}
	return "<!doctype html><html><head><meta charset=\"utf-8\"></head><body>" + html + "</body></html>"
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "document.pdf"
	}

	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "/", "_")

	sanitized := strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return '_'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		switch r {
		case '-', '_', '.':
			return r
		default:
			return '_'
		}
	}, name)

	sanitized = strings.Trim(sanitized, " ._-")
	if sanitized == "" {
		sanitized = "document"
	}
	if len(sanitized) > 120 {
		sanitized = sanitized[:120]
	}
	if !strings.HasSuffix(strings.ToLower(sanitized), ".pdf") {
		sanitized += ".pdf"
	}
	return sanitized
}

func blockedURLPatterns() []string {
	return []string{
		"http://*/*",
		"https://*/*",
		"file://*/*",
		"ftp://*/*",
		"ws://*/*",
		"wss://*/*",
		"chrome://*/*",
		"chrome-extension://*/*",
	}
}
