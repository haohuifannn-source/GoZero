# 通过go-zero框架去实现一个商城的服务

## 当存在多个目录的go.mod时

会出现飘红，无法定位相关的位置，需要通过在最外层用go work init去创建一个文件夹，然后在文件夹里面填入存在的go.mod，注意这里是相对你打开的那个子目录位置的

## 通过datasource去链接数据库

```go
goctl model mysql datasource --url="root:root@tcp(127.0.0.1:3306)/mall" --table="user" --dir="./model" --style="goZero"------如果一直被拒绝，换成下面的方式

goctl model mysql datasource --url="root:root@tcp(172.23.80.1:3306)/mall" --table="user" --dir="./model" --style="goZero"
```

## 用户登录等操作----service\user
(1). 注册操作
参数校验-雪花算法加密-加盐密码-存入