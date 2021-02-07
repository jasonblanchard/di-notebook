ARG BUILD_DIR=/go/src/github.com/jasonblanchard/di-notebook

FROM golang:1.14 AS build
ARG BUILD_DIR

WORKDIR ${BUILD_DIR}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o build/cli -v ./cmd/cli
RUN go build -o build/grpc -v ./cmd/grpc
RUN go build -o build/http -v ./cmd/http
RUN go build -o build/nats -v ./cmd/nats

FROM ubuntu AS run
ARG BUILD_DIR

RUN useradd -ms /bin/bash docker
USER docker

WORKDIR /cmd
ENV PATH="/app:${PATH}"

COPY --from=build --chown=docker:docker ${BUILD_DIR}/build .
COPY --from=build --chown=docker:docker ${BUILD_DIR}/cmd/http/config.yaml .
# TODO: Figure out how/when to get this by tag
COPY --from=build --chown=docker:docker ${BUILD_DIR}/cmd/http/notebook.swagger.json .
COPY --chown=docker:docker migrations migrations/
