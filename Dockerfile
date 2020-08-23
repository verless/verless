# This Dockerfile is designed for hacking on verless - the image
# distributed on Docker Hub is built using `application.Dockerfile`
# in the scripts directory.
#
# Build the image:
#   $ docker image build -t verless .
#
# For hacking on verless, mount the source code against /src:
#   $ docker container run -v $(pwd):/src verless
#
# This will run `go run cmd/main.go`, meaning your code changes
# become visible immediately.
#
# Example: Change config.Version in to "my-version" and run the
# version command by appending it to the image name:
#
#   $ docker container run -v $(pwd):/src verless version
#     verless version my-version
#
# Also, the verless binary is available as `verless` command.
# Note that this is the version from the image build.
FROM golang:1.14-alpine

LABEL maintainer="Dominik Braun <mail@dominikbraun.io>"

RUN mkdir /project
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o ./target/verless cmd/main.go && \
    cp ./target/verless /bin/verless

ENTRYPOINT ["go", "run", "cmd/main.go"]
