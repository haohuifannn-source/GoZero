package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// mysql
	Mysql struct { //数据库配置，出了mysql之外，还可能存在mongo等数据库
		DataSource string //Mysql的链接地址
	}
	// redis
	CacheRedis cache.CacheConf
}
