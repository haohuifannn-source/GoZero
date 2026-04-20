// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	//"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ShortUrlMapModel
	SquenceModel      sequence.Sequence
	SquenceRedis      sequence.Sequence
	ShortUrlBlackList map[string]struct{} //这是为了找黑名单的时候可以不遍历查找，提高效率

	// 手动redis实现
	BizRedis *redis.Redis // 专门处理业务逻辑的 Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	shortConn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	// 把配置文件中配置的黑名单加载到map，方便后续的判断
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	// 手动redis配置
	bizRedis := redis.MustNewRedis(c.BizRedis)

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(shortConn, c.CaCheRedis),
		SquenceModel:      sequence.NewMySql(c.Sequence.DSN, "a"),
		SquenceRedis:      sequence.NewRedis(c.SequenceRedis.Host, "sequence:a:"),
		ShortUrlBlackList: m,
		BizRedis:          bizRedis,
	}

}
