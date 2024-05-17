.PHONY: up down fresh logs test lint migrate sqlc

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

logs:
	docker compose logs -f

test:
	go test -v -race -cover -count=1 -failfast ./...

lint:
	golangci-lint run -v

sqlc:
	docker run --rm -u $$(id -u):$$(id -g) -v ./:/src -w /src sqlc/sqlc $(filter-out $@,$(MAKECMDGOALS))

%:
    @: