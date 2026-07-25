FROM node:24-bookworm AS web-build
WORKDIR /src/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.25-bookworm AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/jianli-server ./cmd/server

FROM debian:bookworm-slim AS app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates chromium fontconfig fonts-liberation fonts-noto-cjk fonts-wqy-microhei tzdata \
    && rm -rf /var/lib/apt/lists/*
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=go-build /out/jianli-server /app/bin/jianli-server
RUN mkdir -p /app/web/dist /app/data/uploads/avatars
COPY --from=web-build /src/web/dist /app/web/dist
# Copy the ip2region xdb if present in the build context. The wildcard
# pattern keeps the build green when the file is absent; the server then
# falls back to the loopback-only resolver and logs a warning.
COPY data/ip2region_v*.xdb* /app/data/
EXPOSE 8080
CMD ["/app/bin/jianli-server"]
