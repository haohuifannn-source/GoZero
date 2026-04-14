// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"time"

	"api/internal/svc"
	"api/internal/types"
	"api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	// 1、 判断传入的用户是否存在
	un := req.Username
	pwd := req.Password
	// 1.1 不存在则返回
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, un)
	if err == sqlx.ErrNotFound {
		logx.Errorw("用户不存在，请注册", logx.Field("err", err))
		return nil, errors.New("用户名或密码错误")
	}
	if err != nil {
		logx.Errorw("Login find user failed", logx.Field("username", un))
		return nil, errors.New("内部错误")
	}
	// 2、 判断该用户的密码是否正确
	if !utils.PasswordVerify(pwd, u.Password) {
		logx.Infow("用户密码错误", logx.Field("username", un))
		return nil, errors.New("用户名或密码错误")
	}

	//3、 生成JWT
	now := time.Now().Unix()
	acessexpir := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := utils.GenerateAccessToken(l.svcCtx.Config.Auth.AccessSecret, now, acessexpir, u.UserId)
	if err != nil {
		logx.Errorw("GenerateAccessToken failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}

	refreshexpir := l.svcCtx.Config.RefreshExpire
	refreshToken, err := utils.GenerateRefreshToken(l.svcCtx.Config.RefreshSecret, now, refreshexpir, u.UserId)
	if err != nil {
		logx.Errorw("GenerateRefreshToken failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}

	return &types.LoginResponse{
		Message:      "登录成功",
		AccessToken:  accessToken,
		ExpiresIn:    int(now + acessexpir),
		RefreshToken: refreshToken}, nil
}
