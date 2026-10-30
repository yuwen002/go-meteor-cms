# ---------- 前端构建阶段 ----------
FROM node:20-alpine AS web
WORKDIR /web
COPY web/package.json ./
RUN npm install
COPY web/ ./
# 产物位于 /web/dist（vite 默认 base 为 /，由 Go 服务统一托管 SPA）
RUN npm run build

# ---------- 后端构建阶段 ----------
FROM golang:1.26-alpine AS builder

WORKDIR /src

# 先复制依赖清单，利用镜像层缓存加速后续构建
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Go 直接从目录托管前端静态文件，无需 embed；dist 由运行阶段挂载
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/cms ./api/cms/v1 \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/migrate ./cmd/migrate

# ---------- 运行阶段 ----------
FROM alpine:3.20

RUN apk add --no-cache tzdata ca-certificates \
 && adduser -D -u 10001 app

WORKDIR /app

COPY --from=builder /out/cms /app/cms
COPY --from=builder /out/migrate /app/migrate
COPY --from=web /web/dist /app/web/dist
# 将 OpenAPI 文档一并发布到静态目录，供 /docs.html（Swagger UI）读取
COPY openapi.json /app/web/dist/openapi.json
COPY api/cms/v1/etc/cms-api.yaml /app/etc/cms-api.yaml

RUN chown -R app:app /app
USER app

ENV TZ=Asia/Shanghai
# 由 Go 服务直接托管的前端产物目录（SPA + history 路由回退）
ENV WEB_DIST=/app/web/dist

EXPOSE 8888

CMD ["/app/cms", "-f", "etc/cms-api.yaml"]
