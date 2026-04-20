package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Redis struct {
	store *redis.Redis
	key   string
}

func NewRedis(redisAdr string, keyStr string) Sequence {
	return &Redis{
		store: redis.New(redisAdr),
		key:   keyStr,
	}
}

func (r *Redis) Next() (seq uint64, err error) {
	val, err := r.store.Incr(r.key)
	if err != nil {
		logx.Errorw("r.store.Incr failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	return uint64(val), nil

}
