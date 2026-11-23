FROM node:24-bookworm AS web-build
WORKDIR /src/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM node:24-bookworm AS geoip-data
COPY deploy/scripts/download-geoip.sh /download-geoip.sh
RUN bash /download-geoip.sh /out/geoip

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
# Keep bundled databases outside /app/data so its host mount cannot hide them.
COPY --from=geoip-data /out/geoip /app/geoip
EXPOSE 8080
CMD ["/app/bin/jianli-server"]
