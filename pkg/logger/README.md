# 日志组件

封装zap作为日志组件

## 使用

```go
logger.Default().Info("info")
logger.With(zap.String("key", "value")).Info("info")

// 把日志实例放入context中
ctx:=logger.ContextWith(ctx, zap.String("key", "value"))
ctx:=logger.ContextWithLogger(ctx, logger.With(zap.String("key", "value")))

// 从context中取出日志实例使用
logger.Ctx(ctx).Info("info")
logger.CtxWith(ctx, zap.String("key", "value")).Info("info")
```