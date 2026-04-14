// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

	"api/internal/svc"
	"api/internal/types"
	"api/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshReq) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	token, err := jwt.ParseWithClaims(req.RefreshToken, jwt.MapClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(l.svcCtx.Config.RefreshSecret), nil
		})
	if err != nil || !token.Valid {
		return nil, errors.New("401 Refresh Token 已失效或非法")
	}
	// 2. 安全提取 userId
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无法解析身份载荷")
	}

	// 建议使用这种方式防止 float64 导致的转换问题
	var userId int64
	if val, ok := claims["userId"]; ok {
		switch v := val.(type) {
		case float64:
			userId = int64(v)
		case string:
			userId, _ = strconv.ParseInt(v, 10, 64)
		}
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		logx.Errorw("Refresh user not found", logx.Field("userId", userId))
		return nil, errors.New("用户状态异常")
	}
	now := time.Now().Unix()
	acessexpir := l.svcCtx.Config.Auth.AccessExpire
	newAccess, _ := utils.GenerateAccessToken(
		l.svcCtx.Config.Auth.AccessSecret, now, acessexpir, user.UserId)
	refreshexpir := l.svcCtx.Config.RefreshExpire
	newRefresh, _ := utils.GenerateRefreshToken(
		l.svcCtx.Config.RefreshSecret, now, refreshexpir, user.UserId)

	return &types.LoginResponse{
		Message:      "刷新成功",
		AccessToken:  newAccess,
		ExpiresIn:    int(now + acessexpir),
		RefreshToken: newRefresh}, nil
}
