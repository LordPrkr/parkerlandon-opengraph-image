# OpenGraph Image Generator

A Go HTTP server that dynamically generates OpenGraph images for blog posts and web pages. The server renders HTML/CSS templates using [Templ](https://templ.guide/) and captures screenshots using a headless Chrome browser via [Rod](https://go-rod.github.io/).

## Features

- Dynamic OG image generation from query parameters
- 1200x630 pixel output (standard OpenGraph dimensions)
- Browser pool for handling concurrent requests
- Embedded static assets for single-binary deployment
- TLS support
- Configurable via TOML

## Quick Start

```bash
# Install dependencies
just deps

# Generate TLS certificates (optional, requires mkcert)
just gen-certs "localhost 127.0.0.1"

# Run the server
just run
```

## Usage

### Generate an OG Image

```bash
# With TLS
curl -k "https://localhost:3000/?title=My%20Blog%20Post&subtitle=A%20great%20article" -o og-image.png

# Without TLS (update config to disable)
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
    │  Generator  │  Get browser from pool
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
│   │   └── pool.go      # Browser pool management
│   ├── server/          # HTTP server setup and routing
│   └── templates/       # Templ HTML templates
├── web/static/          # Embedded static assets (avatar)
└── configs/server/      # TOML configuration files
```

### Key Components

**Browser Pool (`pkg/ogimage/pool.go`)**

- Manages a pool of headless Chrome browser instances
- Channel-based semaphore pattern for concurrency control
- Configurable pool size via `browser_pool_size` config

**Generator (`pkg/ogimage/generator.go`)**

- Acquires browser from pool
- Navigates to internal `/_preview` route
- Captures viewport screenshot at 1200x630

**Template (`pkg/templates/ogimage.templ`)**

- Defines the OG image layout using Templ
- Pink/magenta borders at top and bottom
- Avatar on left (33% width), text content on right
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
browser_pool_size=4                    # number of browser instances
base_url="http://localhost:3000"       # internal URL for screenshots
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

- Go 1.23+
- [just](https://github.com/casey/just) (task runner)
- [mkcert](https://github.com/FiloSottile/mkcert) (optional, for TLS)

### Commands

```bash
# Install Go dependencies and templ CLI
just deps

# Generate Templ templates (creates *_templ.go files)
just generate

# Run the server (includes template generation)
just run

# Generate TLS certificates
just gen-certs "localhost 127.0.0.1"
```

### Customizing the Template

1. Edit `pkg/templates/ogimage.templ`
2. Run `just generate` to regenerate Go code
3. Restart the server

### Replacing the Avatar

Replace `web/static/avatar.png` with your own image and rebuild:

```bash
cp /path/to/your/image.png web/static/avatar.png
just run
```

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
