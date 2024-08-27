FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY ./components/consumers/dependency-track/dependency-track /app/components/consumers/dependency-track/dependency-track

ENTRYPOINT ["/app/components/consumers/dependency-track/dependency-track"]
