ARG base_image=alpine:edge

FROM golang:1 as builder
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go mod download
RUN go build -o /assets/in ./cmd/in
RUN go build -o /assets/out ./cmd/out
RUN go build -o /assets/check ./cmd/check

# there are no tests, but all resources must have a 'tests' target, so just
# no-op
FROM scratch AS tests

FROM ${base_image}
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
RUN apk --no-cache add bash wget netcat-openbsd
