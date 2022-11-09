# build stage
FROM golang:1.18.1-alpine as builder
LABEL "com.conekta.vendor"="Conekta"
LABEL "com.conekta.maintainer"="Franklin Carrero <franklin.carrero@conekta.com>"
LABEL "version"="2021.0.1"


RUN apk update \
  && apk add bash ca-certificates git openssh gcc g++ libc-dev librdkafka-dev pkgconf make curl

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
WORKDIR /go/src/risk-rules
# Copy all the Code and stuff to compile everything
COPY go.mod go.sum ./

RUN --mount=type=ssh git config --global url."ssh://git@github.com/conekta".insteadOf https://github.com/conekta && go mod download -x
# Copy all the Code and stuff to compile everything
COPY . .
# Downloads all the dependencies in advance (could be left out, but it's more clear this way)

RUN \
    # Builds the application as a static linked one, to allow it to run on alpine
    GOOS=linux \
    GOARCH=amd64 \
    go build  -tags musl,appsec  -o compiled-app  ./cmd/httpserver/main

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest
# `service` should be replaced here as well
COPY --from=builder /go/src/risk-rules/compiled-app .

ENV DD_SERVICE="risk-rules" \
    DD_TRACE_ENABLED="true" \
    DEFAULT_PAGE_SIZE="10" \
    DEFAULT_PROCESS_STATUS="1" \
    DD_APPSEC_ENABLED="true" \
    RETRIES="10" \
    TIMEOUT="2s"

CMD ["./compiled-app"]
