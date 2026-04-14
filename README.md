# 学习go-zero框架

## 介绍

api文件 --> goctl工具 --> 一键生成

其中go-zero的github访问口为：https://github.com/zeromicro/go-zero/blob/master/readme-cn.md

## goctl工具

### 安装

可以通过以下命令去安装goctl

```go
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

安装完后通过以下命令去验证安装是否成功

```go
goctl --version
```

### 安装依赖

如果需要生成RPC代码，需要安装protoc和go插件，参考李文周的博客进行安装

简单的方法

```go
goctl env install
```


## 目录组成

（1）gozero_demo是根据官方文档的快速入门

（2）mall是通过一个例子去熟悉api的调用

