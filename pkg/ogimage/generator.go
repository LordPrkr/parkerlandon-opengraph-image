package ogimage

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/ParkerGits/go-backend-starter/pkg/templates"
	"github.com/go-rod/rod/lib/proto"
)

type Generator struct {
	browser *BrowserManager
	baseURL string
}

func NewGenerator(browser *BrowserManager, baseURL string) *Generator {
	return &Generator{
		browser: browser,
		baseURL: baseURL,
	}
}

func (g *Generator) Generate(ctx context.Context, params templates.OGImageParams) ([]byte, error) {
	browser := g.browser.Get()
	if browser == nil {
		return nil, fmt.Errorf("failed to get browser")
	}

	slog.Info("creating page")
	page, err := browser.Timeout(30 * time.Second).Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer page.MustClose()

	slog.Info("setting viewport")
	err = page.Timeout(10 * time.Second).SetViewport(&proto.EmulationSetDeviceMetricsOverride{
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

	slog.Info("navigating to preview", "url", previewURL)
	err = page.Timeout(30 * time.Second).Navigate(previewURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	slog.Info("waiting for page load")
	err = page.Timeout(30 * time.Second).WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for load: %w", err)
	}

	slog.Info("taking screenshot")
	screenshot, err := page.Timeout(30 * time.Second).Screenshot(true, &proto.PageCaptureScreenshot{
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

	slog.Info("screenshot complete")
	return screenshot, nil
}
