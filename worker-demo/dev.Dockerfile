FROM golang:1.24-alpine

RUN go install github.com/air-verse/air@v1.61.7

WORKDIR /app

CMD ["air", "-c", ".air.toml", "--", "run"]
