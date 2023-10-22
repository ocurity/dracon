FROM golang:1.21.3-bookworm as builder
COPY . /build
WORKDIR /build
RUN go build -o components/producers/aggregator/tagger ./components/producers/aggregator/main.go

FROM scratch
COPY --from=builder /build/components/producers/aggregator/tagger /bin/tagger

ENTRYPOINT ["/bin/tagger"]
