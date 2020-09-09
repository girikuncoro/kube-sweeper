FROM golang:1.15-alpine AS build

RUN apk update && \
    apk add build-base

COPY . /build
WORKDIR /build

RUN make build

FROM alpine

COPY --from=build /build/bin/kubesweeper .

ENTRYPOINT ["./kubesweeper"]
