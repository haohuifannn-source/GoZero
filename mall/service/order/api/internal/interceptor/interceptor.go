package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 这种方式避免了硬编码，防止出现重复的定义导致覆盖的情况
type CtxKey string

const (
	CtxKeyAdminID CtxKey = "admin"
)

// MyInterceptor客户端一元拦截器
func MyInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端拦截器启动")
	// RPC调用前
	// 编写客户端拦截器的逻辑
	adminID := ctx.Value(CtxKeyAdminID).(string)
	md := metadata.Pairs(
		"key1", "val1",
		"key1", "val1-2", // "key1"的值将会是 []string{"val1", "val1-2"}
		"requestID", "12345",
		"token", "mall-order-test",
		"userID", adminID,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)           // 将metadata随着RPC发送出去
	err := invoker(ctx, method, req, reply, cc, opts...) // 实际的RPC调用
	// RPC调用后
	fmt.Println("客户端拦截器结束")
	return err
}
