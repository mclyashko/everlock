FROM --platform=linux/amd64 golang:1.23-bullseye AS builder

WORKDIR /everlock

COPY . /everlock/

RUN go mod tidy && go mod vendor

WORKDIR /everlock/cmd/everlock

RUN go build -o ./everlock .

FROM --platform=linux/amd64 debian:bullseye

RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

COPY --from=builder /everlock /everlock

WORKDIR /everlock/cmd/everlock

ENTRYPOINT ["./everlock"]
