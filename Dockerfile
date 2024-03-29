###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build \
      -ldflags='-w -s -extldflags "-static"' \
      -o ./bin/oav \
 && upx-ucl --best --ultra-brute ./bin/oav

###########
# release #
###########

FROM gcr.io/distroless/static-debian11:latest AS release

COPY --from=builder /build/bin/oav /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/oav"]
