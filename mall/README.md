# 通过go-zero框架去实现一个商城的服务

## 1. 当存在多个目录的go.mod时

会出现飘红，无法定位相关的位置，需要通过在最外层用go work init去创建一个文件夹，然后在文件夹里面填入存在的go.mod，注意这里是相对你打开的那个子目录位置的

## 2. 通过datasource去链接数据库

```go
goctl model mysql datasource --url="root:root@tcp(127.0.0.1:3306)/mall" --table="user" --dir="./model" --style="goZero"------如果一直被拒绝，换成下面的方式

goctl model mysql datasource --url="root:root@tcp(172.23.80.1:3306)/mall" --table="user" --dir="./model" --style="goZero"----这里的地址是本机/服务器的地址
```

### 2.1 通过api去生成相应的http文件

```go
goctl api go -api user.api -dir . -style=goZero
```

## 3. 链接云服务器进行开发

(1). 配置文件中的 Host不需要改，保持 0.0.0.0 即可，让你的 user-api 程序监听你本地 WSL 所有的网卡地址
(2). 把mysql、redis这些的网址改为服务器的公网网址
(3). 用postman发请求的时候要用localhost

## 4. 用户登录等操作----service\user

(1). 注册操作
参数校验-雪花算法加密-加盐密码-存入数据库和Redis-集成logx来记录日志

(2). 登录功能

(3). 用户详情功能
这里面可以在user.api里面，用path/form都可以，这里面用的是path

(4). 加入JWT鉴权

用户详情接口需要登录之后才能访问，需要认证auth:

1. 用户成功登录之后

1.1 生成JSON Web Token(JWT) 
1.2 返回给前端，前端代码需要把token保存起来，后续每一次的请求都会带上这个token
2. 后端需要鉴权的接口就会对请求进行鉴权，从请求头中取到token进行解析
2.1. 解析成功就是登录用户
2.2. 解析失败就是未登录或者token失效的用户
3. refresh token（可选）

一般来说生成AccessToken需要带有Role，因为在微服务之中可能会通过RPC进行传输，这时候需要检验该用户是admin还是user，从而在实现delete/Insert操作的时候不需要再次查询数据库。同时，RefreshToken是由前端主动发起，去捕获401来刷新，因此后端只需要实现接口即可。

(5). 加入中间件

1. 路由中间件，即在.api文件中给理由加入中间件
2. 全局中间件，及在middleware文件夹下自定义一个功能文件，如本项目中的global.go

(6). 加入GRPC服务

1. 定义一个.proto文件，定义相关的RPC协议，通过指令生成相关的代码

```go
goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.
```

1. 完善配置结构体Config.go和配置文件yaml(一定要对应上)
2. 完善serviceCtx
3. 完善RPC的逻辑
4. RPC测试工具：grpc ui （[https://github.com/fullstorydev/grpcui）](https://github.com/fullstorydev/grpcui）)

5.1 安装

```go
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

确保环境变量 $GOPATH/bin 目录，添加到环境变量里面

5.2 使用
其中`localhost:12345`是我rpc服务的地址

```
grpcui -plaintext localhost:12345
```

如果出现下面的情况

```bash
haohui@DESKTOP-PC8N4D7:~/projects/goZero/mall/service/user/api$ grpcui -plaintext localhost:8080
Failed to compute set of methods to expose: server does not support the reflection API
```

5.4 如果出现数据库操作失败的时候，需要检查etcd的端口，要配置为以下的形式

```bash
docker run -d --name Etcd-server \
    -p 2379:2379 \
    -p 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://117.72.109.40:2379 \ ---- 这里必须是公网的IP地址
    --env ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
    bitnami/etcd:latest
```

5.3. 使用测试工具的时候，需要把Model定义为dev，因为默认是pro

## 5. 订单服务

(1). 实现一个通过订单的ID调用rpc查询用户信息的功能，通过uerid-->grpc-->user.GetUser

1. 生成model
2. 配置RPC客户端（设置.yaml和.config文件,都需要加入RPC客户端的配置，注意ETCD的key要对应上服务端的key，不是项目的名称）
3. 修改svc文件，告诉代码生成的代码现在又RPC客户端

注意：编写客户端的代码并不需要通过proto文件去生成一个GRPC的通信逻辑，只需要在svc中配置相应的服务即可

## 6. 通过consul实现服务注册和服务的发现

[https://github.com/zeromicro/zero-contrib/tree/main/zrpc/registry](https://github.com/zeromicro/zero-contrib/tree/main/zrpc/registry)

### 服务注册

1. 修改配置文件和yaml，引入consul
  引入"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
   屏蔽ETCD
2. 启动的时候将服务注册到consul
  consul.RegisterService(c.ListenOn, c.Consul)

### 服务发现

1. 修改yaml，引入consul
  Target: consul://117.72.109.40:8500/consul-user.rpc?wait=14s
2. 修改启动代码，引入_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"

## 7. RPC调用传递metadata

### 知识点
1. 什么是metadata？什么样的数据应该存入metadata？它和请求参数有什么区别？

元数据（metadata）是指在处理RPC请求和响应过程中需要但又不属于具体业务（例如身份验证详细信息）的信息，采用键值对列表的形式，其中键是string类型，值通常是[]string类型，但也可以是二进制数据。gRPC中的 metadata 类似于我们在 HTTP headers中的键值对，元数据可以包含认证token、请求标识和监控标签等。
具体可以看https://www.liwenzhou.com/posts/Go/gRPC/#c-6-3

2. GRPC拦截器：客户端的拦截器和服务端的拦截器的区别？

gRPC 为在每个 ClientConn/Server 基础上实现和安装拦截器提供了一些简单的 API。 拦截器拦截每个 RPC 调用的执行。用户可以使用拦截器进行日志记录、身份验证/授权、指标收集以及许多其他可以跨 RPC 共享的功能。

在 gRPC 中，拦截器根据拦截的 RPC 调用类型可以分为两类。第一个是普通拦截器（一元拦截器），它拦截普通RPC 调用。另一个是流拦截器，它处理流式 RPC 调用。而客户端和服务端都有自己的普通拦截器和流拦截器类型。因此，在 gRPC 中总共有四种不同类型的拦截器。
具体可以看https://www.liwenzhou.com/posts/Go/gRPC/#c-6-3

### 客户端拦截器

order 服务的search接口中添加拦截器，添加一些userID、token、requestID等数据

**几个关键点**：
1.什么时候存入metadata

在调用GRPC前需要存入

2.怎么存

通过l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "33")存入

3.拦截器如何通过context传值

通过在拦截器中ctx = metadata.NewOutgoingContext(ctx, md) 将metadata随着RPC发送出去

4.context存取值操作

因为上面通过l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "33")存入了上下文，因此可以通过adminID := ctx.Value(CtxKeyAdminID).(string)取出数据

**注意**：在order的服务里面，用l.ctx = context.WithValue(l.ctx, "adminID", "33")传入一些元数据，这个是给自己的进程看的。和metadata.AppendToOutgoingContext不一样，这个是发到GRPC中给服务端进行调用的

### 服务端拦截器
1.拦截器怎么加？什么时候加？

需要在启动RPC服务之前就添加，s.AddUnaryInterceptors(myInterceptors)

2.拦截器的业务逻辑如何写？

func myInterceptors(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 调用前
	fmt.Println("服务端拦截器启动")

	// 拦截器逻辑
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "need metadata")
	}
	fmt.Printf("metadata %#v\n", md)
	// 根据metadata中的数据进行一些校验处理
	if md["token"][0] != "mall-order-test" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	m, err := handler(ctx, req) // 实际的RPC方法调用

	// 调用后
	fmt.Println("服务端拦截器结束")
	return m, err
}

（1）首先通过metadata.FromIncomingContext(ctx)取出相应的元数据
（2）然后编写相关的拦截逻辑，通过md取出元数据进行验证
（3）然后m, err := handler(ctx, req) 进行实际的RPC方法调用
（4）最后返回

3.服务端拦截器如何从metadata取值

通过md, ok := metadata.FromIncomingContext(ctx)进行取值

## 8. 错误处理
```json
{
	"code": 1001,
	"msg"： "内部错误"
}
```

1. 自定义错误结构体

在errorx的文件夹下面，定义了相关的代码错误和RPC调用错误的结构体

2. 业务逻辑代码中按需返回自定义的错误

在响应返回错误的地方调用自定义的错误函数

3. 告诉go-zero框架处理自定义错误

需要实现效果，就需要在主函数里面调用框架提供的错误处理钩子函数去实现

## 9. 响应的扩展（进阶）

很多情况下返回给前端的数据是以下的格式，带有状态码和返回的msg等，因此需要定义模板
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    ...
  }
}
```

模板的用处是用来生成代码的，是给goctl进行代码生成的时候可以根据模板生成

```go
goctl api go -api user.api -dir . -style=goZero
```

查看默认的存放模板文件的路径
```bash
goctl env
```

通过初始化模板文件得到模板
```bash
goctl template init
```

然后修改，例如下面下面修改handler文件

```go
package {{.PkgName}}

import (
        "net/http"
        "mall/service/order/api/internal/response" // ***1. 加入了自定义的resonse文件***
        "github.com/zeromicro/go-zero/rest/httpx"
        {{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                {{if .HasRequest}}var req types.{{.RequestType}}
                if err := httpx.Parse(r, &req); err != nil {
                        httpx.ErrorCtx(r.Context(), w, err)
                        return
                }

                {{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
                {{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
                {{if .HasResp}}response.Response(w, resp, err){{else}}response.Response(w, nil, err){{end}}//***2. 将原本返回的格式改为自定的***
        }
}
```

然后重新goctl生成一下