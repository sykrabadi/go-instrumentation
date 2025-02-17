# Dockerfile

FROM golang:1.22-alpine AS build
ARG APP_VERSION=latest

WORKDIR /src/

COPY . /src/.

RUN CGO_ENABLED=0 go build -o /bin/go-instrumentation-prometheus ./main.go

# final - service
FROM alpine as api-service

COPY --from=build /bin/go-instrumentation-prometheus /bin/go-instrumentation-prometheus

# HTTP Port
EXPOSE 9320

ENTRYPOINT ["/bin/go-instrumentation-prometheus"]