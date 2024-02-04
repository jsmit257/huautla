# FROM golang:latest as build
# WORKDIR /go/src/build
# RUN git config --global --add safe.directory /go/src/build
# RUN pwd
# CMD CGO_ENABLED=0 go build -v -x -a \
#   -ldflags '-extldflags "-static"' \
#   -o ./internal/cmd/serve-mysql/user-service \
#   ./internal/cmd/serve-mysql/...

FROM debian:bookworm as install
ENTRYPOINT []
VOLUME /huautla
WORKDIR /huautla
RUN apt-get update
RUN apt-get -y upgrade postgresql-client
CMD ./bin/install-prod.sh

FROM debian:bookworm as system-test
VOLUME /system-test
WORKDIR /system-test
CMD ./tests

# FROM alpine:3.14 as deploy
# COPY ./internal/cmd/serve-mysql/user-service /user-service
