ARG PRODUCER_AGGREGATOR_BASE_IMAGE
FROM ${PRODUCER_AGGREGATOR_BASE_IMAGE}

COPY ./components/producers/aggregator /bin/tagger

ENTRYPOINT ["/bin/tagger"]
