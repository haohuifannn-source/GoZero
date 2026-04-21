// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	reidisDb1ShortUrlMapSurlPrefix = "redis:db1:shortUrlMap:surl:"
	Err404                         = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 自己写缓存，就是可以是按surl -->  lurl
// go-zero生成的缓存，就是 surl ---> 数据行，即将模型的结构的数据全部缓存， 会增加缓存的压力

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 查看短链接， 通过短链接重定向到真是的链接
	// 手动添加redis
	// redisKey := fmt.Sprintf("%s%s", reidisDb1ShortUrlMapSurlPrefix, req.ShortURL)
	// s, err := l.svcCtx.BizRedis.Get(redisKey)
	// if err != nil {
	// 	logx.Errorw("l.svcCtx.BizRedis.Get failed", logx.LogField{Key: "err", Value: err.Error()})
	// 	return nil, err
	// }
	// if len(s) > 0 {
	// 	return &types.ShowResponse{LongURL: s}, nil
	// }

	// 1. 根据短链接查询原始的长链接
	// 1.0 布隆过滤器
	// 不存在就返回404，不需要后续的处理
	// a. 基于内存版本，服务重启之后就没了，所以每次的重启都需要去加载一下已有的链接（从数据库查询）
	// b. 基于redis版本，通过go-zero自带的
	// exit, err := l.svcCtx.FilterBloom.Exists([]byte(req.ShortURL))
	// if err != nil {
	// 	logx.Errorw("l.svcCtx.FilterBloom.Exists failed", logx.LogField{Key: "err", Value: err.Error()})
	// 	// 生产建议：如果布隆过滤器报错（比如 Redis 挂了），为了业务可用性，通常选择“放行”而不是报错
	// }

	// // 不存在就返回
	// if !exit {
	// 	return nil, Err404
	// }

	// 布谷鸟过滤器------可选
	// if !l.svcCtx.BooGuFilter.Lookup([]byte([]byte(req.ShortURL))) {
	// 	logx.Errorw("l.svcCtx.BooGuFilter.Lookup failed", logx.LogField{Key: "err", Value: "不存在"})
	// 	return
	// }

	fmt.Println("开始查询缓存和DB......")
	// go-zero本身自带防止缓存击穿的方法
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(
		l.ctx,
		sql.NullString{String: req.ShortURL, Valid: true})
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, Err404
		}
		logx.Errorw("l.svcCtx.ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 手动查到之后还需要去设置
	// _ = l.svcCtx.BizRedis.Setex(redisKey, u.Lurl.String, 86400)

	// 2. 返回查询到的长链接，在调用handler层返回重定向
	return &types.ShowResponse{LongURL: u.Lurl.String}, nil
}
