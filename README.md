# OpenGraph Image Generator

A Go HTTP server that dynamically generates OpenGraph images for blog posts and web pages. The server renders HTML/CSS templates using [Templ](https://templ.guide/) and captures screenshots using a headless Chrome browser via [Rod](https://go-rod.github.io/).

## Features

- Dynamic OG image generation from query parameters
- 1200x630 pixel output (standard OpenGraph dimensions)
- Single shared browser instance for fast page creation
- Embedded static assets for single-binary deployment
- CDN-friendly caching headers
- Configurable via TOML
- Hot-reloading development with Air
- Fly.io deployment ready

## Quick Start

```bash
# Install dependencies
just deps

# Run the server (development)
just dev

# Or run directly
just run
```

## Usage

### Generate an OG Image

```bash
curl "http://localhost:3000/?title=My%20Blog%20Post&subtitle=A%20great%20article" -o og-image.png
```

### Query Parameters

| Parameter  | Required | Description                       |
| ---------- | -------- | --------------------------------- |
| `title`    | Yes      | Main title displayed on the image |
| `subtitle` | No       | Secondary text below the title    |

### Example Response

Successful requests return a PNG image with `Content-Type: image/png`.

Missing required parameters return JSON:

```json
{ "ok": false, "code": -6, "message": "title query parameter is required" }
```

## API Endpoints

| Endpoint    | Method | Description                           |
| ----------- | ------ | ------------------------------------- |
| `/`         | GET    | Generate OG image (returns PNG)       |
| `/_preview` | GET    | HTML preview of the OG image template |
| `/static/*` | GET    | Serve embedded static assets          |

## Architecture

```
Request Flow:

GET /?title=Hello&subtitle=World
         │
         ▼
    ┌─────────────┐
    │   Handler   │  Parse & validate query params
    └─────────────┘
         │
         ▼
    ┌─────────────┐
    │  Generator  │  Get shared browser, create page
    └─────────────┘
         │
         ▼
    ┌─────────────┐
    │     Rod     │  Navigate to /_preview?title=...&subtitle=...
    └─────────────┘
         │
         ▼
    ┌─────────────┐
    │   Templ     │  Render HTML template with CSS
    └─────────────┘
         │
         ▼
    ┌─────────────┐
    │ Screenshot  │  Capture 1200x630 viewport as PNG
    └─────────────┘
         │
         ▼
    Return PNG image
```

### Package Structure

```
├── cmd/server/          # Application entry point
├── pkg/
│   ├── config/          # Viper-based TOML configuration
│   ├── httplib/         # HTTP utilities, middleware, responses
│   ├── ogimage/         # OG image generation logic
│   │   ├── handler.go   # HTTP handler for GET /
│   │   ├── generator.go # Rod screenshot logic
│   │   └── browser.go   # Shared browser management
│   ├── server/          # HTTP server setup and routing
│   └── templates/       # Templ HTML templates
├── web/static/          # Embedded static assets (avatar)
└── configs/server/      # TOML configuration files
```

### Key Components

**Browser Manager (`pkg/ogimage/browser.go`)**

- Manages a single shared headless Chrome browser instance
- Browser launched once at startup (~2-3s)
- Pages created per-request (~100-200ms)
- Auto-reconnects if browser crashes

**Generator (`pkg/ogimage/generator.go`)**

- Creates a new page per request on the shared browser
- Navigates to internal `/_preview` route
- Captures viewport screenshot at 1200x630
- 30-second timeouts on all operations

**Template (`pkg/templates/ogimage.templ`)**

- Defines the OG image layout using Templ
- Pink/magenta borders at top and bottom
- Avatar on left (33% width), text content on right
- Uses system-installed Noto Serif font
- Compiled to Go code via `templ generate`

## Configuration

Configuration files are TOML format in `configs/server/`, selected via the `PROFILE` environment variable (defaults to "local").

### Full Configuration Example

```toml
[cors]
allowed_origins="http://localhost:5173"  # comma-separated for multiple

[tls.server]
cert_file="certs/server.crt"  # omit both to disable TLS
key_file="certs/server.key"

[log]
level="DEBUG"  # DEBUG, INFO, WARN, ERROR

[ogimage]
base_url="http://localhost:3000"       # internal URL for screenshots
browser_path="/usr/bin/chromium"       # path to Chrome/Chromium binary (omit to auto-detect)
```

### Disabling TLS

To run without TLS, remove or comment out the TLS section:

```toml
[tls.server]
# cert_file="certs/server.crt"
# key_file="certs/server.key"
```

## Development

### Prerequisites

- Go 1.22+
- [just](https://github.com/casey/just) (task runner)
- [mkcert](https://github.com/FiloSottile/mkcert) (optional, for TLS)

### Commands

```bash
# Install Go dependencies, templ CLI, and air
just deps

# Run with hot-reloading (recommended for development)
just dev

# Generate Templ templates (creates *_templ.go files)
just generate

# Run the server directly
just run

# Generate TLS certificates
just gen-certs "localhost 127.0.0.1"
```

### Customizing the Template

1. Edit `pkg/templates/ogimage.templ`
2. Air will auto-rebuild, or run `just generate` manually
3. Server restarts automatically

### Replacing the Avatar

Replace `web/static/avatar.png` with your own image and rebuild:

```bash
cp /path/to/your/image.png web/static/avatar.png
just run
```

## Container Deployment

### Build and Run Locally

```bash
# Build container image
just build-image

# Run container
just run-image
```

### Deploy to Fly.io

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login and deploy
fly auth login
fly launch --no-deploy  # first time only
fly deploy
```

The app includes CDN-friendly caching headers:
- Browser cache: 24 hours (`max-age=86400`)
- CDN cache: 7 days (`s-maxage=604800`)
- Stale-while-revalidate: 1 hour

## Output Example

The generated image follows this layout:

```
┌────────────────────────────────────────────────────────┐
│ ██████████████████████████████████████████████████████ │  <- Pink border
│                                                        │
│   ┌─────────┐                                          │
│   │         │      Title Text Here                     │
│   │  Avatar │                                          │
│   │  Image  │      subtitle text here                  │
│   │         │                                          │
│   └─────────┘      @lordprkr                           │
│                                                        │
│ ██████████████████████████████████████████████████████ │  <- Pink border
└────────────────────────────────────────────────────────┘
        1200px x 630px
```

## License

MIT
