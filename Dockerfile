#################
# Base image
#################
FROM alpine:3.22.0 as alicloud-status-rss-base

USER root

RUN addgroup -g 10001 alicloud-status-rss && \
    adduser --disabled-password --system --gecos "" --home "/home/alicloud-status-rss" --shell "/sbin/nologin" --uid 10001 alicloud-status-rss && \
    mkdir -p "/home/alicloud-status-rss" && \
    chown alicloud-status-rss:0 /home/alicloud-status-rss && \
    chmod g=u /home/alicloud-status-rss && \
    chmod g=u /etc/passwd
RUN apk add --update --no-cache alpine-sdk curl

ENV USER=alicloud-status-rss
USER 10001
WORKDIR /home/alicloud-status-rss

#################
# Builder image
#################
FROM golang:1.23-alpine AS alicloud-status-rss-builder
RUN apk add --update --no-cache alpine-sdk
WORKDIR /app
COPY . .
RUN make build

#################
# Final image
#################
FROM alicloud-status-rss-base

COPY --from=alicloud-status-rss-builder /app/bin/alicloud-status-rss /usr/local/bin

# Command to run the executable
ENTRYPOINT ["alicloud-status-rss"]