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

RUN apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    zip \
    unzip \
    build-essential \
    time \
    && rm -rf /var/lib/apt/lists/*

# C/C++(gcc-12)
RUN apt-get update && apt-get install --no-install-recommends -y \
    gcc-12 \
    g++-12 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /judge-server
COPY --from=builder /work/app .

ENTRYPOINT [ "/judge-server/app" ]
