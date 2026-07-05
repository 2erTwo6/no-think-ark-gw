# no-think-ark-gw

一个简单的 API 网关，转发请求到火山引擎 Ark 推理服务，并自动注入 `thinking: {type: "disabled"}` 以禁用深度思考。

## 使用

```bash
podman run -d --restart=always -p 8080:8080 --name ark-gateway ark-gateway
```

请求原样转发到 `https://ark.cn-beijing.volces.com/api/v3/chat/completions`，仅注入 `thinking` 与 `stream` 字段。

## 构建

```bash
podman build -t ark-gateway .
```

## 环境变量

- `PORT` — 监听端口，默认 `8080`
