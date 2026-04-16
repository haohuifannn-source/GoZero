// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"mall/service/order/api/internal/interceptor"
	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 查询数据库得到订单的信息
	// oID, _ := strconv.ParseUint(req.OrderID, 10, 64)
	// o, err := l.svcCtx.OrderMdel.FindOne(l.ctx, oID)
	// if err != nil {
	// 	if err == sqlx.ErrNotFound {
	// 		logx.Errorw("该订单不存在", logx.Field("err", err))
	// 		return nil, errors.New("订单错误")
	// 	}
	// 	logx.Errorw("FindOne failed", logx.Field("err", err))
	// 	return nil, errors.New("内部错误")
	// }
	// 2. 调用rpc服务查询user的信息int64(o.UserId)
	// 在调用RPC前要把metadata存入上下文中
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "33")
	in := &userclient.GetUserReq{UserID: 38255257658068992}
	rpcRsp, err := l.svcCtx.UserRPC.GetUser(l.ctx, in)
	if err != nil {
		logx.Errorw("UserRPC.GetUser failed", logx.Field("err", err))
		return nil, errors.New("远程调用错误")
	}

	return &types.SearchResponse{
		Message:  "查询成功",
		OrderID:  "166666", //strconv.FormatUint(o.OrderId, 10)
		Username: rpcRsp.Username,
		Status:   100, //int(o.Status)
	}, nil
}
