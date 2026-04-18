// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"shorturl/internal/svc"
	"shorturl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShorturlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShorturlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShorturlLogic {
	return &ShorturlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShorturlLogic) Shorturl(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	if req.ShortURL == "1ly7vk" {
		return &types.Response{LongURL: "https://www.liwenzhou.com/posts/Go/golang-menu"}, nil
	}
	return &types.Response{LongURL: "https://www.baidu.com"}, nil
}
