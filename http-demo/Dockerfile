## Multistage build
FROM golang:1.22 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src
COPY . .
RUN  go mod download
RUN go build -o /app

## Multistage deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /src/configs /configs
COPY --from=build /src/templates /templates
COPY --from=build /app /app

ENTRYPOINT ["/app"]
