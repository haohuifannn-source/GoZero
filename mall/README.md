# 通过go-zero框架去实现一个商城的服务

## 当存在多个目录的go.mod时

会出现飘红，无法定位相关的位置，需要通过在最外层用go work init去创建一个文件夹，然后在文件夹里面填入存在的go.mod，注意这里是相对你打开的那个子目录位置的

## 通过datasource去链接数据库

```go
goctl model mysql datasource --url="root:root@tcp(127.0.0.1:3306)/mall" --table="user" --dir="./model" --style="goZero"------如果一直被拒绝，换成下面的方式

goctl model mysql datasource --url="root:root@tcp(172.23.80.1:3306)/mall" --table="user" --dir="./model" --style="goZero"
```

### 通过api去生成相应的http文件
```go
goctl api go -api user.api -dir . -style=goZero
```

## 用户登录等操作----service\user
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