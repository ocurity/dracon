ARG BASE_IMAGE
FROM ${BASE_IMAGE:-scratch}

ARG GOOS
ARG GOARCH

COPY ./components/producers/aggregator/${GOOS}/${GOARCH}/aggregator-parser /app/components/producers/aggregator/tagger

ENTRYPOINT ["/app/components/producers/aggregator/tagger"]
