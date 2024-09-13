FROM alpine:3.20

ARG GOOS
ARG GOARCH

RUN apk add --no-cache ca-certificates

COPY ./components/producers/github-code-scanning/${GOOS}/${GOARCH}/github-code-scanning-parser /app/components/producers/github-code-scanning/github-code-scanning-parser

ENTRYPOINT ["/app/components/producers/github-code-scanning/github-code-scanning-parser"]
