
ARG golang_version

#FROM golang:$golang_version

FROM golang:1.9.0-alpine3.6

RUN echo "Build number: $golang_version"

MAINTAINER Alexey Kovrizhkin <lekovr+docker@gmail.com>

RUN apk add --no-cache make bash git g++ curl

WORKDIR /go/src/github.com/LeKovr/teleproxy
# Will fetch git commit ID
COPY .git .git
# Sources
COPY cmd cmd
# make build
COPY Makefile .

#sqlite3 is a cgo package
#ENV CGO_ENABLED=0

ENV GOOS=linux
ENV BUILD_FLAG=-a
#"-tags netgo -a -v"

RUN make build

# ------------------------------------------------------------------------------

# Cant use it because sqlite3
#FROM scratch
FROM alpine:3.6

RUN apk add --no-cache make bash curl

WORKDIR /
COPY --from=0 /go/src/github.com/LeKovr/teleproxy/cmd/teleproxy/teleproxy .
# Need for SSL
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Templates sample
COPY messages.gohtml /

ENTRYPOINT ["/teleproxy"]

