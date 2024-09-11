FROM alpine:3.20

RUN apk add --no-cache ca-certificates

ARG GOOS=linux
ARG GOARCH=amd64

COPY ./components/consumers/slack/${GOOS}/${GOARCH}/slack /app/components/consumers/slack/slack

ENTRYPOINT ["/app/components/consumers/slack/slack"]
