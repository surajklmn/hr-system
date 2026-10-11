FROM golang:alpine AS builder
WORKDIR /app

# copy go.mod and go.sum and download dependencies first
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build without C dependency and name the compiled executable server and then rn run
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

# goest to step 1 takes the compiled executables and copies it, source code and caches aren't copied
COPY --from=builder /app/server .

# the port we want to use
EXPOSE 8080

# this command runs when the container boots up
CMD ["./server"]

