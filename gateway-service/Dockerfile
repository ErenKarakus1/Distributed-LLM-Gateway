FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/server

FROM alpine:3.22

WORKDIR /app
COPY --from=build /app/gateway .

EXPOSE 8080

CMD ["./gateway"]
