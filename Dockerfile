# ---------- 构建阶段 ----------
FROM golang:1.26-alpine AS builder

WORKDIR /src

# 先复制依赖清单，利用镜像层缓存加速后续构建
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/cms ./api/cms/v1 \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/migrate ./cmd/migrate

# ---------- 运行阶段 ----------
FROM alpine:3.20

RUN apk add --no-cache tzdata ca-certificates \
 && adduser -D -u 10001 app

WORKDIR /app

COPY --from=builder /out/cms /app/cms
COPY --from=builder /out/migrate /app/migrate
COPY api/cms/v1/etc/cms-api.yaml /app/etc/cms-api.yaml

RUN chown -R app:app /app
USER app

ENV TZ=Asia/Shanghai

EXPOSE 8888

CMD ["/app/cms", "-f", "etc/cms-api.yaml"]
