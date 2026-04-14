// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"errors"
	"mall/service/user/model"

	"api/internal/svc"
	"api/internal/types"
	"api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	// todo: add your logic here and delete this line
	//参数校验---应该放在controller层里面
	if req.Password != req.RePassword {
		return nil, errors.New("两次输入的密码不一致")
	}
	logx.Infov(req) //json.Marshall(req)
	logx.Infof("req:%#v\n", req)
	// 把用户的注册信息保存到数据库中
	// 0、 查询username是否已经被注册
	//https://github.com/go-sql-driver/mysql一些关于数据库链接的参数
	username := req.Username
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw(
			"FindOneByUsername failed",
			logx.Field("err", err),
		)
		return nil, errors.New("内部错误")
	}
	if u != nil {
		return nil, errors.New("用户名已经存在")
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Gender:   int64(req.Gender),
	}
	// 1、 生成userID（雪花算法）---后续需要放在protoc里面做
	uuId := utils.GenerateUserID()
	user.UserId = uuId

	// 2、加密密码（加盐|md5）
	newpassword, err := utils.PasswordHash(user.Password)
	if err != nil {
		logx.Errorw(
			"PasswordHash failed",
			logx.Field("err", err),
		)
		return nil, err
	}
	user.Password = newpassword

	if _, err := l.svcCtx.UserModel.Insert(context.Background(), user); err != nil {
		logx.Errorw(
			"Insert failed",
			logx.Field("err", err),
		)
		return nil, err
	}
	return &types.SignupResponse{Message: "sucess"}, nil
}
