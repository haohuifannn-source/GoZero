package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// 定义一个全局中间件

// 功能：
// 记录所有请求的响应信息

// rest.Middle -> type Middleware func(next http.HandlerFunc) http.HandlerFunc
// type HandlerFunc func(ResponseWriter, *Request)

type bodyCopy struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func NewBodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (bc *bodyCopy) Write(b []byte) (int, error) {
	// 1. 先在我的小本本记录响应
	bc.body.Write(b)
	// 2. 往HTTP响应里面记录响应内容
	return bc.ResponseWriter.Write(b)
}

// CopyResp复制请求的响应体
func CopyResp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 处理请求前,得到自定义的http.ResponseWriter
		bc := NewBodyCopy(w)

		next(bc, r) // 实际的路由处理handler函数
		// 处理请求后
		fmt.Printf("----> req : %v resp : %v\n", r.URL, bc.body.String())
	}
}

// 调用其他服务的中间件
func MiddlewareWithAnotherService(ok bool) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 处理请求前,得到自定义的http.ResponseWriter
			if ok {
				fmt.Println("----->ok!")
			}
			next(w, r)
		}
	}
}
