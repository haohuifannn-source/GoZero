// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserDetailByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserDetailByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailByIdLogic {
	return &GetUserDetailByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserDetailByIdLogic) GetUserDetailById(req *types.UserIDRequest) (resp *types.UserMsgResponse, err error) {
	// todo: add your logic here and delete this line
	// 1、判断userID是否存在
	uId := req.Userid
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, uId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Infow("userID not exit", logx.Field("userID", uId))
			return nil, errors.New("内部错误")
		}
		logx.Errorw("UserModel.FindOneByUserId failed", logx.Field("userID", uId))
		return nil, errors.New("内部错误")
	}
	// 2、返回结果
	return &types.UserMsgResponse{Username: u.Username, Gender: int(u.Gender)}, nil
}
