FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o notification ./cmd

EXPOSE 8082

CMD ["./notification"]
