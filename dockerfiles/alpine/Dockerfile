FROM golang:1 as builder
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go mod download
RUN go build -o /assets/in ./cmd/in
RUN go build -o /assets/out ./cmd/out
RUN go build -o /assets/check ./cmd/check

FROM alpine:edge
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
RUN apk --no-cache add bash wget
