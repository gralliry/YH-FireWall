# Web 服务器

## 前端构建

```bash
# 编译前端到 dist/
cd ../../frontend
npm run build

# 将产物复制到 embed 目录
cp -r dist/* ../internal/server/webserver/static/

# 编译 Go 二进制（产物在 build/ 下）
cd ../..
go build -o build/yfwd ./cmd/yfwd
```

前端源码位于 `frontend/`，编译产物输出到 `frontend/dist/`，然后复制到 `internal/server/webserver/static/` 被 Go embed 打包。

## Swagger 文档

```bash
swag init -g cmd/yfwd/main.go --parseDependency --parseDepth 2 -o internal/server/webserver/docs
```

Swagger 文档生成在 `internal/server/webserver/docs/`，通过 `GET /swagger/*` 访问。
