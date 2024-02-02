ARG PRODUCER_AGGREGATOR_BASE_IMAGE
FROM ${PRODUCER_AGGREGATOR_BASE_IMAGE}

COPY ./components/producers/aggregator/aggregator-parser /app/components/producers/aggregator/tagger

ENTRYPOINT ["/app/components/producers/aggregator/tagger"]
