[TOC]

# Cago

Cago 一个快速开发的集成式框架.使用模块化的开发模式,每一个组件都可以单独的调用.

Cago 只对社区工具进行集成,大大减少迁移难度和学习成本,我们不生产代码,我们只是方案的搬运工.

使用 go 的`struct`来声明 API 和 swagger 文档,可以通过脚手架来帮助你生成相关内容,大大减轻 API 开发的困难.

## 快速开始

[简单示例](./examples/simple)

## 脚手架

[脚手架文档](./cmd/cago)

## CI/CD

[部署资源](./deploy)

cago 提供了[`gitlab-ci`](deploy/gitlab/.gitlab-ci.yml)、[`golanglint-ci`](./deploy/.golangci.yml)和
[`Kubernetes`](./deploy)的 CI/CD 配置文件,可以快速帮你实现 CI/CD.

当本地调试时也可以使用`docker compose up -d`启动调试环境.

默认使用`etcd`作为配置中心,同时也支持文件作为配置启动.

## 组件

- [常用组件包](./pkg/component)
- [数据库组件包](./database)
- [服务组件包](./server)
- [中间件](./middleware)
- [工具包/杂项](./pkg)

## 参考项目

- GoFrame
- GoMicro

## License

[MIT](./LICENSE)
