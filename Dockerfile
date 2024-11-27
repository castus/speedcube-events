FROM golang:1.23.3 AS builder

WORKDIR /data
COPY . /data
RUN GOOS=linux GOARCH=amd64 go build -o speedcube-events

FROM ubuntu:noble
RUN apt update -y && apt upgrade -y && apt install -y \
  locales \
  curl \
  tzdata \
  ca-certificates \
  openssl \
  nano
RUN apt clean && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY --from=builder /data/speedcube-events ./
