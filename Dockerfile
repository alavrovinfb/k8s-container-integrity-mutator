FROM golang:1.18.5-alpine AS build
ENV CGO_ENABLED 0

WORKDIR /src
ADD . ./
RUN go mod download
RUN go build -o k8s-webhook-injector ./cmd

RUN chmod +x k8s-webhook-injector

FROM alpine:latest
WORKDIR /app
COPY --from=build /src/k8s-webhook-injector .

ENTRYPOINT ["/app/k8s-webhook-injector"]
