# This Dockerfile builds a lightweight distribution image to be
# released on Docker Hub. It contains the built application, but
# no source code.
FROM alpine:3.11.5 AS downloader

# The verless version to download.
ARG VERSION

RUN apk add --no-cache \
    curl \
    tar

RUN curl -LO https://github.com/verless/verless/releases/download/${VERSION}/verless-linux-amd64.tar && \
    tar -xzvf verless-linux-amd64.tar -C /bin && \
    rm -f verless-linux-amd64.tar

# The final stage which corresponds to the distributed image.
FROM alpine:3.11.5 AS final

COPY --from=downloader ["/bin/verless", "/bin/verless"]

WORKDIR /project

ENTRYPOINT ["/bin/verless"]
CMD ["build", "/project"]