FROM concourse/golang-builder as builder
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go mod download
RUN go build -o /assets/in ./cmd/in
RUN go build -o /assets/out ./cmd/out
RUN go build -o /assets/check ./cmd/check

FROM ubuntu:bionic AS resource
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
RUN apt-get update && apt-get -y install wget netcat-openbsd
