FROM golang:1.21.3 AS builder

WORKDIR /data
COPY . /data
RUN GOOS=linux GOARCH=amd64 go build -o speedcube-events

FROM ubuntu:mantic
COPY --from=builder /data/speedcube-events ./
CMD ["./speedcubing-calendar"]
