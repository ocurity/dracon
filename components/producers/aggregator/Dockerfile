ARG PRODUCER_AGGREGATOR_BASE_IMAGE
FROM ${PRODUCER_AGGREGATOR_BASE_IMAGE}

ARG GOOS=linux
ARG GOARCH=amd64

COPY ./components/producers/aggregator/${GOOS}/${GOARCH}/aggregator-parser /app/components/producers/aggregator/tagger

ENTRYPOINT ["/app/components/producers/aggregator/tagger"]
