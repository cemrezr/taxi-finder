FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/driver-location-api

RUN go build -o main .

WORKDIR /app/cmd/matching-api

RUN go build -o main .


CMD ["./main"]
