# Cago 示例

这是一个简单的示例，提供了一个目录结构与组件使用的最佳实践，包含了以下组件的使用：

- logger 日志
- mysql 数据库
- nsq 消息队列

在运行前你需要使用docker将开发环境启动起来：`docker compose up -d`

## 运行

请注意你的运行目录，需要将[`config.yaml`](./configs/config.yaml.tmp)复制到运行目录下的`configs`文件夹中
