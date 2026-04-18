// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"shortener/internal/config"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ShortUrlMapModel
	SquenceModel  model.SequenceModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	shortConn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	squenceConn := sqlx.NewMysql(c.Sequence.DSN)
	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewShortUrlMapModel(shortConn),
		SquenceModel:  model.NewSequenceModel(squenceConn),
	}
}
