ARG BASE_IMAGE
FROM ${BASE_IMAGE:-scratch}

COPY ./components/producers/aggregator/aggregator-parser /app/components/producers/aggregator/tagger

ENTRYPOINT ["/app/components/producers/aggregator/tagger"]
