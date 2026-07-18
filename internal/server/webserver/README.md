# Web 服务器

## 前端构建

```bash
# 编译前端到 dist/
cd ../../ui
npm run build

# 编译 Go 二进制（产物在 build/ 下，自动 embed ui/dist/）
cd ../..
go build -o build/yfwd ./cmd/yfwd
```

前端源码位于 `ui/`，编译产物输出到 `ui/dist/`。Go 编译时通过 `ui/embed.go` 的 `//go:embed dist/*` 自动嵌入，无需手动复制。

## Swagger 文档

```bash
swag init -g cmd/yfwd/main.go --parseDependency --parseDepth 2 -o internal/server/webserver/docs
```

Swagger 文档生成在 `internal/server/webserver/docs/`，通过 `GET /swagger/*` 访问。
