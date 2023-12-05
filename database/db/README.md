# 数据库

底层使用gorm库进行封装，支持多种数据库、支持单库与多库模式。

~~多库模式基于`gorm.io/plugin/dbresolver`实现。~~

## 配置

```yaml
# 单库模式, key设置为db
driver: mysql
dsn: root:password@tcp(127.0.0.1:3306)/db?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local&multiStatements=true
prefix: prefix_

# 多库模式, key设置为dbs
default: # 默认链接, 必须设置
  driver: mysql
  dsn: root:password@tcp(127.0.0.1:3306)/db?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local&multiStatements=true
  prefix: prefix_
ch: # clickhouse
  driver: clickhouse
  dsn: clickhouse://127.0.0.1:9009/default?read_timeout=10s

```

## 使用

```go
db.Default().Model(&User{}).Where("id = ?", 1).First(&user)
db.Ctx(ctx).Model(&User{}).Where("id = ?", 1).First(&user)

// 多库
db.With("clickhouse").Model(&User{}).Where("id = ?", 1).First(&user)
db.CtxWith(ctx, "clickhouse").Model(&User{}).Where("id = ?", 1).First(&user)
```

## 驱动

默认支持`mysql`，其它驱动需要使用`db.RegisterDriver`进行注册。可以参考[clickhouse](./clickhouse.go)的实现。

```go
import (
_ "github.com/codfrm/cago/database/db/clickhouse"
)
```
