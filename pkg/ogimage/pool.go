package ogimage

import (
	"log/slog"
	"os"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// BrowserPool manages browser instances. In constrained environments,
// it creates browsers on-demand rather than pre-pooling them.
type BrowserPool struct {
	mu sync.Mutex
}

func NewBrowserPool(size int) (*BrowserPool, error) {
	// Size is ignored - we create browsers on-demand
	return &BrowserPool{}, nil
}

// Get creates a fresh browser for each request
func (p *BrowserPool) Get() *rod.Browser {
	p.mu.Lock()
	defer p.mu.Unlock()

	browser, err := launchBrowser()
	if err != nil {
		slog.Error("failed to launch browser", "error", err)
		return nil
	}
	return browser
}

// Put closes the browser after use
func (p *BrowserPool) Put(browser *rod.Browser) {
	if browser != nil {
		browser.MustClose()
	}
}

func (p *BrowserPool) Close() {
	// Nothing to close - browsers are closed after each request
}

func launchBrowser() (*rod.Browser, error) {
	if browserPath := os.Getenv("ROD_BROWSER"); browserPath != "" {
		slog.Info("launching browser", "path", browserPath)
		u, err := launcher.New().
			Bin(browserPath).
			Headless(true).
			NoSandbox(true).
			Set("disable-gpu").
			Set("disable-dev-shm-usage").
			Set("disable-software-rasterizer").
			Set("single-process").
			Launch()
		if err != nil {
			return nil, err
		}
		return rod.New().ControlURL(u).MustConnect(), nil
	}

	return rod.New().MustConnect(), nil
}
