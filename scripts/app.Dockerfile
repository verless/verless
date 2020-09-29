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
    tar -xvf verless-linux-amd64.tar -C /bin && \
    rm -f verless-linux-amd64.tar

# The final stage which corresponds to the distributed image.
FROM alpine:3.11.5 AS final

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="verless"
LABEL org.label-schema.description="A simple and lightweight Static Site Generator."
LABEL org.label-schema.url="https://github.com/verless/verless"
LABEL org.label-schema.vcs-url="https://github.com/verless/verless"
LABEL org.label-schema.version=${VERSION}
LABEL org.label-schema.docker.cmd="docker container run -v $(pwd)/my-blog:/project verless/verless"

COPY --from=downloader ["/bin/verless", "/bin/verless"]

# Create a symlink for musl, see https://stackoverflow.com/a/35613430.
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /project

ENTRYPOINT ["/bin/verless"]
CMD ["build", "/project"]