// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string //JWT密钥
		AccessExpire int64  //有效期
	}

	RefreshSecret string
	RefreshExpire int64

	Mysql struct { //数据库配置，出了mysql之外，还可能存在mongo等数据库
		DataSource string //Mysql的链接地址
	}

	CacheRedis cache.CacheConf
}
