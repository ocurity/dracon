FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY ./components/consumers/slack/slack /app/components/consumers/slack/slack

ENTRYPOINT ["/app/components/consumers/slack/slack"]