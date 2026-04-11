// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"
	"mall/service/user/model"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 把用户的注册信息保存到数据库中
	// 0、 查询username是否已经被注册
	// 1、 生成userID（雪花算法）
	// 2、加密密码（加盐|md5）
	fmt.Printf("req : %#v\n", req)
	user := &model.User{
		Id:       1111,
		Username: req.Username,
		Password: req.Password,
		Gender:   1,
	}
	if _, err := l.svcCtx.UserModel.Insert(context.Background(), user); err != nil {
		fmt.Printf("Insert failed : %v\n", err)
		return nil, err
	}
	return &types.SignupResponse{Message: "sucess"}, nil
}
