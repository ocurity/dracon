FROM alpine:3.20

RUN apk add --no-cache ca-certificates

ARG GOOS=linux
ARG GOARCH=amd64

COPY ./components/consumers/jira/${GOOS}/${GOARCH}/jira /app/components/consumers/jira/jira

ENTRYPOINT ["/app/components/consumers/jira/jira"]
