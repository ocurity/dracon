FROM alpine:3.20

RUN apk add --no-cache ca-certificates

ARG GOOS=linux
ARG GOARCH=amd64

COPY ./components/consumers/dependency-track/${GOOS}/${GOARCH}/dependency-track /app/components/consumers/dependency-track/dependency-track

ENTRYPOINT ["/app/components/consumers/dependency-track/dependency-track"]
