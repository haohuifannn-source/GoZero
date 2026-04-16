package main

import (
	"context"
	"flag"
	"fmt"

	"mall/service/user/rpc/internal/config"
	"mall/service/user/rpc/internal/server"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 将服务注册到consul
	consul.RegisterService(c.ListenOn, c.Consul)

	defer s.Stop()

	// 注册服务端拦截器
	s.AddUnaryInterceptors(myInterceptors)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start() // 启动RPC服务
}

func myInterceptors(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 调用前
	fmt.Println("服务端拦截器启动")

	// 拦截器逻辑
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "need metadata") // status是GRPC中专门返回错误的方式，带有GRPC状态码，不能用errors.New()
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
