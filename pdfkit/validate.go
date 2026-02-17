package pdfkit

import (
	"fmt"
	"strings"
)

func ValidateRequest(req RenderRequest) error {
	if strings.TrimSpace(req.HTML) == "" {
		return fmt.Errorf("html is required")
	}
	if req.Size != "" {
		key := strings.ToUpper(strings.TrimSpace(req.Size))
		if _, ok := SupportedSizes[key]; !ok {
			return fmt.Errorf("unsupported size: %s", req.Size)
		}
	}
	if req.Orientation != "" {
		key := strings.ToLower(strings.TrimSpace(req.Orientation))
		if _, ok := SupportedOrientations[key]; !ok {
			return fmt.Errorf("unsupported orientation: %s", req.Orientation)
		}
	}
	return nil
}
