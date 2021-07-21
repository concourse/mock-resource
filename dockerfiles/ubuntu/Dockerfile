ARG base_image
ARG builder_image=concourse/golang-builder

FROM ${builder_image} as builder
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
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
RUN apt update && apt upgrade -y -o Dpkg::Options::="--force-confdef"
RUN apt update && apt -y install wget netcat-openbsd
