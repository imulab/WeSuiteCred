FROM golang:1.20.10-alpine3.18 AS builder

ARG BUILD_VERSION
ARG BUILD_REVISION
ARG BUILD_TIME

ENV HOME /build
ENV GOOS linux

WORKDIR $HOME

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build \
    -a \
    -ldflags "-w -s -X absurdlab.io/WeSuiteCred/buildinfo.Version=${BUILD_VERSION} -X absurdlab.io/WeSuiteCred/buildinfo.Revision=${BUILD_REVISION} -X absurdlab.io/WeSuiteCred/buildinfo.CompiledAt=${BUILD_TIME}" \
    -tags urfave_cli_no_docs \
    -o WeSuiteCred \
    .

FROM alpine:3.18 AS runtime

RUN apk add --no-cache curl

COPY --from=builder /build/WeSuiteCred /usr/bin/WeSuiteCred

CMD ["WeSuiteCred", "listener"]
