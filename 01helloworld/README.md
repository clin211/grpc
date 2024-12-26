# hello world grpc
>
> 如果还没有安装对应的环境，要先安装环境，相关文档如下
>
> - [protocol buffer compiler](https://grpc.io/docs/protoc-installation/)
> - [Go Generated Code Guide](https://protobuf.dev/reference/go/go-generated/)

## 克隆项目

```sh
git clone https://github.com/clin211/grpc
```

### 运行项目

因为这个项目有多个语言实现，所以要进入对应的目录运行和查看源码！

#### Go

进入 go 目录，安装依赖，然后运行，如果修改 protoc 后，可以使用命令 `make go-protoc` 来重新生成 pb 文件；生成后的文件在 `/go/rpc` 目录下。

> 本项目的环境：
>
> - 开发系统：MacOS
> - protoc 29.2
> - go 1.22.0
> - vs code 1.96.2
> - node 20.10.0
> - pnpm 8.12.1

- 进入 go 目录

  ```sh
  cd go
  ```

- 安装依赖

  ```sh
  go mod tidy
  ```

- 运行服务端

  ```sh
  go run server/main.go
  ```

- 运行客户端

  ```sh
  go run client/main.go
  ```

- 修改 protoc 后，重新生成 pb 文件

  ```sh
  make go-protoc
  ```

#### node

首先要进入 node 目录：

```sh
cd node
```
>
> 因为使用了 `@grpc/proto-loader` 库，它可以根据最新的 proto 文件自动更新 pb 文件，所以只需要关注业务代码就 OK！如果使用静态生成 pb 文件的话，`grpc-tools` 工具对 ESM 的支持不好，如果是 CommonJS 规范则没什么问题。建议使用项目中的 loader 来自动更新。Nest.js 也是在 `@grpc/proto-loader` 基础上构建的。

- 安装依赖

  ```sh
  pnpm  i
  ```

- 运行服务端

  ```sh
  pnpm server
  ```

- 运行客户端

  ```sh
  pnpm client
  ```

  运行后终端打印结果如下：

  ```txt
  Greeting: Hello clina
  Greeting: Hello again, clina
  ```
