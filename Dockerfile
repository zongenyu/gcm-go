#FROM ubuntu:18.04
FROM alpine:3.8

MAINTAINER ray

RUN mkdir -p /home/eyeon-user/servers/event/gcm-go
WORKDIR /home/eyeon-user/servers/event/gcm-go
ADD ["gcm-go","gcmKeys.json", "cserver", "./"]


RUN apk add --no-cache libc6-compat rsyslog curl \
 && rm -rf /var/cache/apk/*

# CMD ["./cserver"]