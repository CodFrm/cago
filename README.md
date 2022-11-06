# CaGO

## 核心组件包

- [config](./configs)
- [registry](./registry.go)
- [http](./server/http)
- [logger](./pkg/logger)

## 常用组件包

- [broker](./pkg/broker)
- [trace](./pkg/trace)
- [utils](./pkg/utils)

## 数据库

- [sql db](./database/db)

## 快速启动

使用goland打开项目,复制configs/config.yaml.example到configs/config.yaml,修改配置文件

启动example/simple/main.go,即可运行一个简单的服务

另外使用`docker-compose up -d`可以启动框架相关服务(loki、jaeger、grafana、etcd、etcdkeeper)

## 生成工具

在`internal/api`目录下,定义好api请求结构,使用下面命令和自动生成`controller`代码和`swagger`文档

在`internal/service`目录下,定义好service接口,使用下面命令和自动生成`service`代码

```bash
cago gen
```

定义好表结构和`configs`文件后,使用下面命令和自动生成`model`代码和`repository`接口

```bash
cago gen table_name
```

# 参考项目

- GoFrame
- GoMicro

# License

[MIT](LICENSE)