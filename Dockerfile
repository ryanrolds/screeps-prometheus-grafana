FROM golang:1.18-alpine

WORKDIR /collector

COPY . .

RUN go build -o ./bin/scraper ./cmd/scraper

CMD ["./bin/scraper"]