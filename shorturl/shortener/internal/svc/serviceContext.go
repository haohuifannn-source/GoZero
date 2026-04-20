// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ShortUrlMapModel
	SquenceModel      sequence.Sequence
	SquenceRedis      sequence.Sequence
	ShortUrlBlackList map[string]struct{} //这是为了找黑名单的时候可以不遍历查找，提高效率
}

func NewServiceContext(c config.Config) *ServiceContext {

	shortConn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	// 把配置文件中配置的黑名单加载到map， 方便后续的判断
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(shortConn),
		SquenceModel:      sequence.NewMySql(c.Sequence.DSN, "a"),
		SquenceRedis:      sequence.NewRedis(c.SequenceRedis.Host, "sequence:a:"),
		ShortUrlBlackList: m,
	}

}
