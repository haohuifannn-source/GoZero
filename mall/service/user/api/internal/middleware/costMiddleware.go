// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type CostMiddleware struct {
}

func NewCostMiddleware() *CostMiddleware {
	return &CostMiddleware{}
}

func (m *CostMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		now := time.Now()
		// Passthrough to next handler if need
		next(w, r)
		fmt.Printf("----> Cost : %v\n", time.Since(now))
	}
}
