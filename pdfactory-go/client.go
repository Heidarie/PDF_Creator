package pdfactory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout        = 20 * time.Second
	defaultMaxResponseBty = int64(25 << 20)
)

type Client struct {
	baseURL          string
	httpClient       *http.Client
	userAgent        string
	maxResponseBytes int64
}

type Option func(*Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.httpClient.Timeout = timeout
		}
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = strings.TrimSpace(ua)
	}
}

func WithMaxResponseBytes(max int64) Option {
	return func(c *Client) {
		if max > 0 {
			c.maxResponseBytes = max
		}
	}
}

func NewClient(baseURL string, opts ...Option) (*Client, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return nil, errors.New("baseURL is required")
	}
	baseURL = strings.TrimRight(baseURL, "/")

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		maxResponseBytes: defaultMaxResponseBty,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) Render(ctx context.Context, req RenderRequest) ([]byte, error) {
	if strings.TrimSpace(req.HTML) == "" {
		return nil, errors.New("html is required")
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	url := c.baseURL + "/render"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if c.userAgent != "" {
		httpReq.Header.Set("User-Agent", c.userAgent)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	limit := c.maxResponseBytes
	if limit <= 0 {
		limit = defaultMaxResponseBty
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr ErrorResponse
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("render failed: %s", apiErr.Error)
		}
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("render failed: %s", msg)
	}

	if len(body) == 0 {
		return nil, errors.New("empty pdf response")
	}

	return body, nil
}

func (c *Client) RenderToFile(ctx context.Context, req RenderRequest, path string) error {
	pdf, err := c.Render(ctx, req)
	if err != nil {
		return err
	}
	return writeFile(path, pdf)
}
