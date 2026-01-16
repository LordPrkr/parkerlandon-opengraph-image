package ogimage

import (
	"log/slog"
	"os"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// BrowserPool manages a single shared browser instance.
// Pages are created per-request, which is fast (~100ms).
type BrowserPool struct {
	browser *rod.Browser
	mu      sync.Mutex
}

func NewBrowserPool(size int) (*BrowserPool, error) {
	browser, err := launchBrowser()
	if err != nil {
		return nil, err
	}
	slog.Info("browser launched successfully")
	return &BrowserPool{browser: browser}, nil
}

// Get returns the shared browser instance
func (p *BrowserPool) Get() *rod.Browser {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if browser is still alive, reconnect if needed
	if p.browser == nil {
		slog.Info("browser was nil, relaunching")
		browser, err := launchBrowser()
		if err != nil {
			slog.Error("failed to relaunch browser", "error", err)
			return nil
		}
		p.browser = browser
	}

	return p.browser
}

func (p *BrowserPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.browser != nil {
		p.browser.MustClose()
		p.browser = nil
	}
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
			Launch()
		if err != nil {
			return nil, err
		}
		return rod.New().ControlURL(u).MustConnect(), nil
	}

	return rod.New().MustConnect(), nil
}
