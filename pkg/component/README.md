# Cago 组件包

> Cago 组件包,提供框架常用的一些组件

## Core

`component.Core`,核心组件包,提供了框架所需核心组件的初始化

- logger 日志组件,使用zap进行封装
- trace 链路追踪,支持jaeger和uptrace
- metrics 指标监控

## Database

`component.Database`,GORM数据库组件包

- 使用gorm进行封装,支持常见sql数据库

## Mongo

`component.Mongo`,MongoDB数据库组件包

## Redis

`component.Redis`,Redis组件包

## Cache

`component.Cache`,缓存组件包

支持下面的缓存

- redis

## Broker

`component.Broker`,消息队列组件包

支持下面的消息队列

- nsq
- event_bus

