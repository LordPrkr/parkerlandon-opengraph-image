package ogimage

import (
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type BrowserPool struct {
	browsers chan *rod.Browser
	size     int
}

func NewBrowserPool(size int) (*BrowserPool, error) {
	pool := &BrowserPool{
		browsers: make(chan *rod.Browser, size),
		size:     size,
	}

	for range size {
		browser, err := launchBrowser()
		if err != nil {
			return nil, err
		}
		pool.browsers <- browser
	}

	return pool, nil
}

func launchBrowser() (*rod.Browser, error) {
	// Check for ROD_BROWSER env var for system-installed browser
	if browserPath := os.Getenv("ROD_BROWSER"); browserPath != "" {
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

	// Fall back to auto-download behavior
	return rod.New().MustConnect(), nil
}

func (p *BrowserPool) Get() *rod.Browser {
	return <-p.browsers
}

func (p *BrowserPool) Put(browser *rod.Browser) {
	p.browsers <- browser
}

func (p *BrowserPool) Close() {
	close(p.browsers)
	for browser := range p.browsers {
		browser.MustClose()
	}
}
