FROM golang:alpine AS build-env
COPY . /work
WORKDIR /work
RUN go build

FROM alpine
COPY --from=build-env /work/channel_monitor /sxbin/channel_monitor
WORKDIR /sxbin
