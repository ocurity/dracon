ARG OSSF_SCORECARD_SAFETY_BASE_IMAGE
FROM gcr.io/openssf/scorecard:stable

FROM ${OSSF_SCORECARD_SAFETY_BASE_IMAGE}

COPY --from=0 /scorecard /scorecard

ENTRYPOINT ["/scorecard"]