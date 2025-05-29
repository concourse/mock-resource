ARG base_image=cgr.dev/chainguard/wolfi-base
ARG builder_image=concourse/golang-builder

FROM ${builder_image} AS builder

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

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

FROM ${base_image} AS resource
RUN apk --no-cache add bash
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
