FROM golang:1.20-alpine

RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

CMD ["air", "-c", ".air.toml", "--", "run"]
