package ogimage

import (
	"log/slog"
	"sync"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// BrowserManager manages a single shared browser instance.
// Pages are created per-request, which is fast (~100ms).
type BrowserManager struct {
	browser     *rod.Browser
	browserPath string
	mu          sync.Mutex
}

func NewBrowserManager(browserPath string) (*BrowserManager, error) {
	browser, err := launchBrowser(browserPath)
	if err != nil {
		return nil, err
	}
	slog.Info("browser launched successfully")
	return &BrowserManager{browser: browser, browserPath: browserPath}, nil
}

// Get returns the shared browser instance, relaunching if needed
func (m *BrowserManager) Get() *rod.Browser {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if browser is still alive by attempting a simple operation
	if m.browser != nil {
		_, err := m.browser.Version()
		if err != nil {
			slog.Warn("browser connection lost, relaunching", "error", err)
			m.browser = nil
		}
	}

	// Relaunch if needed
	if m.browser == nil {
		slog.Info("launching browser")
		browser, err := launchBrowser(m.browserPath)
		if err != nil {
			slog.Error("failed to launch browser", "error", err)
			return nil
		}
		m.browser = browser
	}

	return m.browser
}

func (m *BrowserManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.browser != nil {
		m.browser.MustClose()
		m.browser = nil
	}
}

func launchBrowser(browserPath string) (*rod.Browser, error) {
	if browserPath != "" {
		slog.Info("launching browser", "path", browserPath)
		u, err := launcher.New().
			Bin(browserPath).
			Headless(true).
			Set("no-sandbox").
			Set("disable-gpu").
			Set("disable-dev-shm-usage").
			Set("disable-setuid-sandbox").
			Set("disable-extensions").
			Set("disable-background-networking").
			Set("disable-sync").
			Set("disable-translate").
			Set("mute-audio").
			Set("hide-scrollbars").
			Set("metrics-recording-only").
			Set("no-first-run").
			Launch()
		if err != nil {
			return nil, err
		}
		return rod.New().ControlURL(u).MustConnect(), nil
	}

	return rod.New().MustConnect(), nil
}
