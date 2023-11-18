FROM gcr.io/distroless/static:nonroot
COPY edge-receiver /
USER 65532:65532

ENTRYPOINT ["/edge-receiver"]
