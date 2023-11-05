FROM golang:1.21.3-alpine as builder

ARG PKG_NAME=github.com/JayJamieson/striphtml

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

COPY . /go/src/${PKG_NAME}

RUN cd /go/src/${PKG_BASE}/${PKG_NAME} && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /striphtml ./cmd/striphtml/

FROM scratch

COPY --from=builder /striphtml .
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
# COPY --from=builder /etc/ssl/certs/ca-bundle.trust.crt /etc/ssl/certs/ca-bundle.trust.crt

ENV PORT=8080

ARG VERSION="local"
ENV ENV_VERSION=$VERSION

EXPOSE 8080

ENTRYPOINT [ "./striphtml", "serve" ]
