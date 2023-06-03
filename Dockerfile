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

COPY --from=busybox /bin/sleep /bin/
COPY --from=busybox /bin/printenv /bin/
COPY --from=busybox /bin/env /bin/
COPY --from=busybox /bin/mkdir /bin/
COPY --from=busybox /bin/nslookup /bin/
COPY --from=busybox /bin/touch /bin/
COPY --from=busybox /bin/true /bin/
COPY --from=busybox /bin/false /bin/
COPY --from=busybox /bin/find /bin/
COPY --from=busybox /bin/mkfifo /bin/
COPY --from=busybox /bin/sed /bin/
COPY --from=busybox /bin/wc /bin/

COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*
