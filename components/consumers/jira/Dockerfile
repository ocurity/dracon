FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY ./components/consumers/jira/jira /app/components/consumers/jira/jira

ENTRYPOINT ["/app/components/consumers/jira/jira"]
