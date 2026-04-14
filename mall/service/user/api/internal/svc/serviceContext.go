// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config    config.Config
	Cost      rest.Middleware //自定义路由中间件，其中字段名要和.api文件中的声明一致
	UserModel model.UserModel //加入增删改查模型
}

func NewServiceContext(c config.Config) *ServiceContext {
	// UserModel -> 接口类型
	// *defaultUserModel 实现了接口
	// 调用构造函数得到model *defaultUserModel，调用New方法

	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlxConn, c.CacheRedis),
		Cost:      middleware.NewCostMiddleware().Handle,
	}
}
