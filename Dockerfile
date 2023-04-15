# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM} golang:1.19-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apk add build-base

WORKDIR /src

COPY . /src

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /backpack cmd/main.go

FROM --platform=${TARGETPLATFORM} alpine:latest as runner

COPY --from=builder /backpack ./

RUN sh -c "apk add sqlite && mkdir /data && \
    mkdir /config && touch /config/rclone.conf && \
    mkdir -p /root/.config/rclone && \
    ln -s /config/rclone.conf /root/.config/rclone/rclone.conf"

CMD ["./backpack", "-config", "/config/config.json"]