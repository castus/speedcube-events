FROM golang:1.22.0 AS builder

WORKDIR /data
COPY . /data
RUN GOOS=linux GOARCH=amd64 go build -o speedcube-events

FROM ubuntu:mantic
RUN apt-get update -y && apt-get upgrade -y && apt-get install -yq \
  locales \
  curl \
  tzdata \
  ca-certificates \
  openssl \
  nano
RUN apt-get clean && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY --from=builder /data/speedcube-events ./
CMD ["./speedcubing-calendar"]
