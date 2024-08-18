FROM golang:alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o stark8 ./cmd/

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/stark8 .
COPY static ./static
EXPOSE 8080
CMD ["./stark8"]
