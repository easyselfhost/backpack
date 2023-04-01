# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM} golang:1.19-alpine as builder

RUN apk add build-base

WORKDIR /src

COPY . /src

RUN sh -c "go get && go build -o /backpack cmd/main.go"

FROM --platform=${BUILDPLATFORM} alpine:latest as runner

COPY --from=builder /backpack ./

RUN sh -c "apk add sqlite && mkdir /data && mkdir /config"

ENV RCLONE_CONFIG=/config/rclone.conf

CMD ["./backpack", "-try-first", "-config", "/config/config.json"]