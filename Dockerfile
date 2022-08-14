FROM golang:1.18-alpine

WORKDIR /collector

COPY . .

RUN go build -o ./bin/screeps-collector ./cmd/screeps-collector

CMD ["./bin/screeps-collector"]