package ogimage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ParkerGits/go-backend-starter/pkg/templates"
	"github.com/go-rod/rod/lib/proto"
)

type Generator struct {
	pool    *BrowserPool
	baseURL string
}

func NewGenerator(pool *BrowserPool, baseURL string) *Generator {
	return &Generator{
		pool:    pool,
		baseURL: baseURL,
	}
}

func (g *Generator) Generate(ctx context.Context, params templates.OGImageParams) ([]byte, error) {
	browser := g.pool.Get()
	defer g.pool.Put(browser)

	page, err := browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer page.MustClose()

	err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  1200,
		Height: 630,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	previewURL := fmt.Sprintf("%s/_preview?title=%s&subtitle=%s",
		g.baseURL,
		url.QueryEscape(params.Title),
		url.QueryEscape(params.Subtitle))

	err = page.Navigate(previewURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	err = page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for load: %w", err)
	}

	screenshot, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  1200,
			Height: 630,
			Scale:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return screenshot, nil
}
