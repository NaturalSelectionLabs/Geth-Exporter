FROM rss3/go-builder as base

WORKDIR /app

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

COPY . .

FROM base AS builder

ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build main.go

FROM rss3/go-runtime AS runner


COPY --from=builder /app/main /bin/geth-exporter

EXPOSE 8000
USER nobody
ENTRYPOINT ["/bin/geth-exporter"]