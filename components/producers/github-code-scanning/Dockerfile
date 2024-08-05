FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY ./components/producers/github-code-scanning/github-code-scanning-parser /app/components/producers/github-code-scanning/github-code-scanning-parser

ENTRYPOINT ["/app/components/producers/github-code-scanning/github-code-scanning-parser"]
