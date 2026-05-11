# Hand-built fallback image. The primary distribution channel is `ko`
# (see .ko.yaml + release.yml). Use this Dockerfile for local builds:
#   docker build -t stackit-nuke:dev .

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=0.0.0-dev
ARG COMMIT=dirty
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w \
      -X github.com/qaiser42/stackit-nuke/pkg/common.SUMMARY=${VERSION} \
      -X github.com/qaiser42/stackit-nuke/pkg/common.COMMIT=${COMMIT}" \
    -o /out/stackit-nuke ./

FROM cgr.dev/chainguard/static:latest
COPY --from=build /out/stackit-nuke /usr/local/bin/stackit-nuke
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/stackit-nuke"]
