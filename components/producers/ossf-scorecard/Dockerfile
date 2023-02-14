FROM gcr.io/openssf/scorecard:stable

FROM alpine:latest

COPY --from=0 /scorecard /scorecard

ENTRYPOINT ["/scorecard"]