FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o sso ./cmd/sso/main.go
RUN go build -o migrator ./cmd/migrator/main.go

CMD ["./sso"]