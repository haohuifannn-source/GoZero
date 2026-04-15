// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"mall/service/order/api/internal/config"
	"mall/service/order/model"
	"mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	OrderMdel model.OrderModel
	UserRPC   userclient.User //是user文件夹下面rpc的userclient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		OrderMdel: model.NewOrderModel(conn, c.CacheRedis),
		UserRPC:   userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
	}
}
