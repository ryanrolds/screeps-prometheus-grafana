FROM golang:1.18-alpine

WORKDIR /scraper

COPY . .

RUN go build -o ./bin/scraper ./cmd/scraper

CMD ["./bin/scraper"]