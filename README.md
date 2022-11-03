# CaGO

## 核心组件包

- [config](configs)
- [registry](./registry)

## 常用组件包

- [logger](./pkg/logger)
- [utils](./pkg/utils)

## 快速启动

使用goland打开项目,复制configs/config.yaml.example到configs/config.yaml,修改配置文件

启动example/simple/main.go,即可运行一个简单的服务

另外使用`docker-compose up -d`可以启动框架相关服务(loki、jaeger、grafana、etcd、etcdkeeper)

# 参考项目

- GoFrame
- go-micro

# License

[MIT](LICENSE)