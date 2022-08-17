FROM golang:1.18-alpine

WORKDIR /scraper

COPY . .

RUN go build -o ./bin/scraper ./cmd/scraper

# TODO - use a multi-stage build to reduce final image size

CMD ["./bin/scraper"]