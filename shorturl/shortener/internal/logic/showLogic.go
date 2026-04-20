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
	redisKey := fmt.Sprintf("%s%s", reidisDb1ShortUrlMapSurlPrefix, req.ShortURL)
	s, err := l.svcCtx.BizRedis.Get(redisKey)
	if err != nil {
		logx.Errorw("l.svcCtx.BizRedis.Get failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	if len(s) > 0 {
		return &types.ShowResponse{LongURL: s}, nil
	}

	// 1. 根据短链接查询原始的长链接
	// go-zero本身自带防止缓存击穿的方法
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(
		l.ctx,
		sql.NullString{String: req.ShortURL, Valid: true})
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, errors.New("404")
		}
		logx.Errorw("l.svcCtx.ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 手动查到之后还需要去设置
	_ = l.svcCtx.BizRedis.Setex(redisKey, u.Lurl.String, 86400)

	// 2. 返回查询到的长链接，在调用handler层返回重定向
	return &types.ShowResponse{LongURL: u.Lurl.String}, nil
}

// 防止缓存击穿的手动实现版
// 缓存击穿一定是发生在缓存查不到的情况下，因为一开始肯定是先要去查，再做缓存击穿防御
// func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
//     // 1. 构造 Redis Key
//     redisKey := fmt.Sprintf("%s%s", reidisDb1ShortUrlMapSurlPrefix, req.ShortURL)

//     // 2. 第一层防御：先查手动 Redis 缓存
//     s, _ := l.svcCtx.BizRedis.Get(redisKey)
//     if len(s) > 0 {
//         return &types.ShowResponse{LongURL: s}, nil
//     }

//     // 3. 第二层防御：使用 singleflight 防止缓存击穿
//     // 这里使用 req.ShortURL 作为 Key，确保相同的短链接只会在同一时间查一次数据库
//     val, err, _ := l.svcCtx.SingleGroup.Do(req.ShortURL, func() (any, error) {
//         // --- 以下逻辑在同一时间只会有一个协程进入 ---

//         // 再次检查缓存（可选，双重检查锁定模式，能进一步提高严谨性）
//         s, _ := l.svcCtx.BizRedis.Get(redisKey)
//         if len(s) > 0 {
//             return s, nil
//         }

//         // 查数据库
//         u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(
//             l.ctx,
//             sql.NullString{String: req.ShortURL, Valid: true},
//         )
//         if err != nil {
//             return nil, err
//         }

//         // 查到之后立刻回填到手动缓存 Redis 中
//         // 这样后面排队的协程被唤醒后，可能直接从缓存拿，或者直接拿这次的结果
//         _ = l.svcCtx.BizRedis.Setex(redisKey, u.Lurl.String, 86400)

//         return u.Lurl.String, nil
//         // --- 以上逻辑结束 ---
//     })

//     if err != nil {
//         if err == sqlx.ErrNotFound {
//             return nil, errors.New("404")
//         }
//         return nil, err
//     }

//     // val 是 any 类型，需要断言回 string
//     return &types.ShowResponse{LongURL: val.(string)}, nil
// }
