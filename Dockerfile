ARG base_image
ARG builder_image=concourse/golang-builder

FROM busybox:uclibc as busybox

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
USER root
COPY --from=builder assets/ /opt/resource/
COPY --from=busybox /bin/* /bin
RUN chmod +x /opt/resource/*
