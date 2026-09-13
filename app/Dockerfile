FROM golang:1.27.1-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/http-server-projeto-korp \
    ./cmd/server

FROM debian:13.6-slim AS runtime

RUN groupadd --system --gid 10001 korp \
    && useradd \
        --system \
        --uid 10001 \
        --gid korp \
        --no-create-home \
        --home-dir /nonexistent \
        --shell /usr/sbin/nologin \
        korp

WORKDIR /app

COPY --from=builder --chown=korp:korp /out/http-server-projeto-korp ./

USER 10001:10001

EXPOSE 8080

ENTRYPOINT ["./http-server-projeto-korp"]