FROM golang:1.19 AS builder

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux \
    GOARCH=amd64 \
    go build \
    -o app \
    .

FROM ubuntu:22.04

WORKDIR /judge-server
COPY --from=builder /work/app .

ENTRYPOINT [ "/judge-server/app" ]
