FROM golang:1.20-alpine

RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

ENTRYPOINT ["air", "-c", ".air.toml", "--"]
CMD ["run"]
