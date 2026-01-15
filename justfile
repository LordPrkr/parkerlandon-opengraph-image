# Generate templ templates
generate:
	templ generate

# Run the server (with generation)
run: generate
	go run cmd/server/main.go

# Generate certs for HOSTS
gen-certs HOSTS:
	rm -rf ./certs
	mkdir ./certs
	mkcert -key-file ./certs/server.key -cert-file ./certs/server.crt {{HOSTS}}

# Install dependencies
deps:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest
	go mod download

# Run development server with hot-reloading 
dev:
	air -c .air.toml

build-image:
	podman build -t opengraph-image .

run-image:
	podman run -it --rm \
		-p 3000:3000 \
		opengraph-image


