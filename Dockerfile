ARG BUILD_DIR=/go/src/github.com/jasonblanchard/di-notebook

FROM golang:1.14 AS build
ARG BUILD_DIR

WORKDIR ${BUILD_DIR}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o build/db -v ./cmd/db

FROM ubuntu AS run
ARG BUILD_DIR

RUN useradd -ms /bin/bash docker
USER docker

WORKDIR /cmd
ENV PATH="/app:${PATH}"

COPY --from=build --chown=docker:docker ${BUILD_DIR}/build .

# CMD ["./nats"]
