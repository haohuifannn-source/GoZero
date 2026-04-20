// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	ShortUrlDB ShortUrlDB

	Sequence struct {
		DSN string
	}

	SequenceRedis SequenceRedis

	BaseString string //指定base64的顺序

	ShortUrlBlackList []string // 指定不能出现的路径词黑名单

	ShortDomain string
}

type ShortUrlDB struct {
	DSN string
}

type SequenceRedis struct {
	Host string
}
