up:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose up -d

down:
	docker compose down

fresh:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V

stubs:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/example.proto
logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...

lint:
	golangci-lint run -v
