FROM golang:1.23 AS BUILDER

WORKDIR /builder
COPY ./ ./
RUN go mod download
RUN ./build.sh

FROM ubuntu:22.04

WORKDIR /srv

COPY --from=BUILDER /builder/platform /srv/

ENTRYPOINT [ "./platform" ]
