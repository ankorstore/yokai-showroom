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

rename:
	find . -type f ! -path "./.git/*" ! -path "./build/*" ! -path "./Makefile" -exec sed -i.bak -e "s|github.com/ankorstore/yokai-http-template|github.com/$(to)|g" {} \;
	find . -type f -name "*.bak" -delete
