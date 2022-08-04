FROM golang:rc-alpine AS build
ENV CGO_ENABLED 0

RUN apk add git openssl

WORKDIR /src
ADD . /src
RUN go mod download
RUN go build -o k8s-webhook-injector cmd/main.go

RUN chmod +x k8s-webhook-injector

FROM scratch
WORKDIR /app
COPY --from=build /src/k8s-webhook-injector .
COPY --from=build /src/certificates/ssl ssl
COPY --from=build /src/patch-json-command.json .

EXPOSE 8443

CMD ["/app/k8s-webhook-injector"]
