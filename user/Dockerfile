FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user ./cmd

EXPOSE 8081

CMD ["./user"]
