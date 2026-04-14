package logic

import (
	"context"
	"errors"

	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	// todo: add your logic here and delete this line
	// 根据userID去查询用户信息
	uID := in.UserID
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, uID)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errors.New("无效的userID")
	}
	if err != nil {
		logx.Errorw("FindOneByUserId failed", logx.Field("err", err))
		return nil, errors.New("查询失败")
	}
	return &user.GetUserResp{
		UserID:   u.UserId,
		Username: u.Username,
		Gender:   u.Gender}, nil
}
