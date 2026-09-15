FROM golang:1.23-alpine AS build

ARG SERVICE_NAME

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY shared ./shared
COPY ${SERVICE_NAME} ./service

WORKDIR /src/service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20

ARG PORT=8080
ENV PORT=${PORT}

RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=build /out/server ./server
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app

USER appuser
EXPOSE ${PORT}

HEALTHCHECK --interval=15s --timeout=3s --start-period=20s --retries=5 \
    CMD wget -qO- "http://127.0.0.1:${PORT}/health" || exit 1

CMD ["./server"]
