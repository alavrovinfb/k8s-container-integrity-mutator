FROM golang:rc-alpine AS build
# need to turn off CGO since otherwise there might be dynamic links
ENV CGO_ENABLED 0

RUN apk add git openssl

WORKDIR /usr/local/go/src/k8s-webhook-injector
ADD certificates .
RUN go mod download
RUN go build ./cmd/k8s-certificates-injector

FROM scratch
WORKDIR /app
COPY --from=build /usr/local/go/src/k8s-webhook-injector/k8s-webhook-injector .
COPY --from=build /usr/local/go/src/k8s-webhook-injector/certificates/ssl ssl
COPY --from=build /usr/local/go/src/k8s-webhook-injector/patch-json-command.json .

EXPOSE 8443

CMD ["/app/k8s-webhook-injector"]
