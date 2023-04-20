FROM golang:1.20.3-alpine3.17 as Builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o sentry-tunnel ./main.go

FROM alpine:latest
ENV APP_HOST=""
ENV APP_PORT="3333"
ENV SENTRY_HOST="localhost"
ENV SENTRY_SCHEMA="https"
ENV SENTRY_PROJECT_IDS=""
ENV APP_ROUTE_PATH="/bugs"
EXPOSE 3333

COPY --from=Builder /app/sentry-tunnel /usr/local/bin/sentry-tunnel
ENTRYPOINT ["sentry-tunnel"]