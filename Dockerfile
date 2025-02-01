FROM golang:1.23 AS builder

WORKDIR /builder
COPY ./ ./
RUN go mod download && \
	go build .

FROM ubuntu:22.04

WORKDIR /srv

COPY --from=builder /builder/platform /srv/
COPY --from=builder /builder/db/schema.sql /srv/db/schema.sql
COPY --from=builder /builder/db/triggers.sql /srv/db/triggers.sql
COPY --from=builder /builder/db/statements.sql /srv/db/statements.sql
COPY --from=builder /builder/templates /srv/templates/
COPY --from=builder /builder/static /srv/static/

ENTRYPOINT [ "/srv/platform" ]
