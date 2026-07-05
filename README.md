# no-think-ark-gw

一个简单的 API 网关，转发请求到火山引擎 Ark 推理服务，并自动注入 `thinking: {type: "disabled"}` 以禁用深度思考。

目的是为了用其免费提供的模型，接入沉浸式翻译。翻译需要高速响应，所以需要关闭thinking。

也许你会问为什么不用New API?也有参数覆写功能啊，但大量的翻译请求会淹没我New API的日志，因此就有了这么一个简单的小项目。

支持高并发，不支持流式传输

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
