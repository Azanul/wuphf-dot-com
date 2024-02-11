FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

COPY . .

RUN go build -o api-gateway ./cmd

EXPOSE 8080

CMD ["./api-gateway"]
