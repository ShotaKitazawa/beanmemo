FROM node:24.14-alpine AS frontend-builder

WORKDIR /app

RUN corepack enable pnpm

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm-store,target=/tmp/pnpm-store \
    pnpm install --frozen-lockfile --store-dir /tmp/pnpm-store

COPY frontend/ .

RUN pnpm run build

FROM golang:1.25 AS backend-builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY backend/ .
COPY --from=frontend-builder /app/dist ./internal/ui/dist

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=backend-builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
