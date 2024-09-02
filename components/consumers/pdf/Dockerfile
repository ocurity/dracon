
# Stage 2: Build
FROM golang:1.22 as builder

# Install playwright cli with right version for later use
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

# Stage 3: Final
FROM ubuntu:oracular
COPY --from=builder /go/bin/playwright /
RUN apt-get update && apt-get install -y ca-certificates tzdata \
    # Install dependencies and all browsers (or specify one)
    && /playwright install chromium --with-deps \
    && rm -rf /var/lib/apt/lists/*

COPY components/consumers/pdf/pdf /application/pdf
COPY components/consumers/pdf/default.html /application/default.html
WORKDIR /application
ENTRYPOINT ["/application/pdf"]
