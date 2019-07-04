# todo micro service

* [[Tutorial, Part 1] How to develop Go gRPC microservice with HTTP/REST endpoint, middleware, Kubernetes deployment, etc.](https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9)

> you should create database first, just run `todo.sql` in db

## 坑

grpc_gateway 中 Unannotated 的方案存在问题

### `*.yaml` 文件配置 http 路由初始化顺序不受控制

比如有两个 API 具有相同的 path pattern

* API1： /v1/todo/all 用于获取所有信息
* API2： /v1/todo/{id} 用于获取指定 id 的信息

路由顺序必须让 API1 先于 API2, 否则会出现字符串 "all" 和 id 类型不一致的问题。

这种情况必须用 annotations 的方式，且 *.proto 中要调整顺序。

### 不支持自定义 openapi 配置

proto 中导入 `protoc-gen-swagger/options/annotations.proto` 可以自定义 openapi 
的一些信息，而 Unannotated 方式未实现这个。

如何自定义 openapi 参见[这个例子](https://github.com/grpc-ecosystem/grpc-gateway/blob/ab0345bb32/examples/proto/examplepb/a_bit_of_everything.proto)
