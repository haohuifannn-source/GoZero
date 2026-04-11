# 通过快速开始来熟悉go-zero框架（更齐全的需要观看官方的文档）

## 第一个greet服务

通过直接生成

```go
goctl api new greet
cd greet
go mod tidy
```

## 第二个user服务

借助api生成

```go
goctl api go -api user.api -dir .
go mod init user-api
go mod tidy
```

## 第三个RPC服务

首先要安装proto的插件，借助代码生成一个服务框架

```go
goctl env install ----如果没有安装protoc插件就必须执行
goctl rpc new user
cd user
go mod tidy
```

然后去修改其中的.proto文件，然后借助下面指令去生成新的文件

```go
goctl rpc protoc user.proto \
    --go_out=./user \
    --go-grpc_out=./user \
    --zrpc_out=./client
```

## 第四个通过自定义的sql语句生成Model

```sql
CREATE TABLE `user` (
  `id`         bigint NOT NULL AUTO_INCREMENT,
  `username`   varchar(255) NOT NULL DEFAULT '',
  `password`   varchar(255) NOT NULL DEFAULT '',
  `mobile`     varchar(20)  NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_mobile`   (`mobile`)
) ENGINE=InnoDB;
```

通过以下的语句去生成响应的Model
```go
goctl model mysql ddl -src user.sql -dir ./internal/model
goctl model mysql ddl -src user.sql -dir ./internal/model -cache
```

