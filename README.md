[TOC]

# CaGO

cago 一个快速开发的集成式框架.使用模块化的开发模式,每一个组件都可以单独的调用.

cago 只对社区工具进行集成,大大减少迁移难度和学习成本,我们不生产代码,我们只是方案的搬运工.

使用 go 的`struct`来声明 API 和 swagger 文档,可以通过脚手架来帮助你生成相关内容,大大减轻 API 开发的困难.

## 核心组件包

> 如果你想使用 cago 来启动你的服务,那么你必须注册以下核心组件

- [config](./configs)
- [registry](./registry.go)
- [http](./server/mux)
- [logger](./pkg/logger)

## 常用组件包

- [broker](./pkg/broker)
- [trace](./pkg/trace)
- [utils](./pkg/utils)

## 数据库

- [sql db](./database/db)
- [redis](./database/redis)

## 中间件

- [session](./middleware/sessions)

# 快速开始

[简单示例](./examples/simple/main.go)

使用 goland 打开项目,复制 configs/config.yaml.example 到 configs/config.yaml,修改配置文件

启动 example/simple/main.go,即可运行一个简单的服务

另外使用`docker-compose up -d`可以启动框架相关服务(loki、jaeger、grafana、etcd、etcdkeeper)

# 脚手架

```bash
go install github.com/codfrm/cago/cmd/cago@latest
```

在`internal/api`目录下,定义好 api 请求结构,使用下面命令和自动生成`controller`代码和`swagger`文档

在`internal/service`目录下,定义好 service 接口,使用下面命令和自动生成`service`代码

```bash
cago gen
```

定义好表结构和`configs`文件后,使用下面命令和自动生成`model`代码和`repository`接口

```bash
cago gen table_name
```

# 部署

cago 提供了`gitlab-ci`、`golanglint-ci`和`Kubernetes`的 CI/CD 配置文件,可以快速帮你实现 CI/CD.

当本地调试时也可以使用`docker-compose up -d`启动调试环境.

默认使用`etcd`作为配置中心,同时也支持文件作为配置启动.

# 参考项目

- GoFrame
- GoMicro

## License

[MIT](LICENSE)
