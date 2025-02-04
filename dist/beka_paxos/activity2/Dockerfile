FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
