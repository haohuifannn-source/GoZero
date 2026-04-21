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
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert转链：输入一个长链接--->转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 校验输入的数据
	// 1.1 数据不可以为空 / 长链接能请求通的网址 / 判断数据库中是否已经存在长链接 / 输入的不能是一个短链接（防止循环转链）

	//1.1.2 长链接能请求通的网址
	if ok := connect.Get(req.LongURL); !ok {
		return nil, errors.New("无效的链接")
	}

	//1.1.3 判断数据库中是否已经存在长链接，生成长链接的MD5值去查数据库
	md5Value := md5.Sum([]byte(req.LongURL))
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, fmt.Errorf("该链接已经被转换为短链接%s", u.Surl.String)
		}
		logx.Errorw("FindOneByMd5 failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}

	// 1.1.4 输入的不能是一个短链接（防止循环转链）
	// 输入的是一个完整的url q1mi.cn/1d12a?name=q1mi
	basePath, err := urltool.GetBasePath(req.LongURL)
	if err != nil {
		logx.Errorw("GetBasePath failed", logx.LogField{Key: "lurl", Value: req.LongURL}, logx.LogField{Key: "err", Value: err})
		return nil, errors.New("内部错误")
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("该链接已经是短链接")
		}
		logx.Errorw("FindOneBySurl failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}

	var short string
	// 为了一直防止生成敏感的路径名称
	for {
		// 2. 取号 ---基于MYSQL实现的发号器
		// 每来一个转链请求，我们就使用replace into 语句往 sequence 表插入一条数据，并取出其中的主键作为号码
		// 基于MySQL实现
		seq, err := l.svcCtx.SquenceModel.Next()
		if err != nil {
			logx.Errorw("l.svcCtx.SquenceModel.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, errors.New("内部错误")
		}
		fmt.Printf("------>Mysqlreq : %v\n", seq)

		// 基于Redis实现
		// redisSeq, err := l.svcCtx.SquenceRedis.Next()
		// if err != nil {
		// 	logx.Errorw("l.svcCtx.SquenceRedis.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
		// 	return nil, errors.New("内部错误")
		// }
		// fmt.Printf("------>Redisreq : %v\n", redisSeq)

		// 3. 号码转短链
		// 3.1 安全性  1En = 6347， 别人可以一直遍历，把你的逻辑扒出来，就是可以知道你的0-62用什么表示
		// 3.2 短域名黑名单避免某些特殊词，比如api\health\fuck
		short = base62.Int2String(seq)
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单的短链接直接跳出for循环
		}

	}
	fmt.Printf("-----> short : %v\n", short)

	// 4. 存储长链接和短链接的映射关系
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongURL, Valid: true},
			Md5:  sql.NullString{String: md5Value, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		},
	); err != nil {
		logx.Errorw("l.svcCtx.ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 将生成的短链接加载到布隆过滤器中
	if err := l.svcCtx.FilterBloom.Add([]byte(short)); err != nil {
		logx.Errorw("l.svcCtx.FilterBloom.Add failed", logx.LogField{Key: "err", Value: err.Error()})
		//生产建议：如果布隆过滤器报错（比如 Redis 挂了），为了业务可用性，通常选择“放行”而不是报错
	}

	// 布谷鸟过滤器
	// if err := l.svcCtx.BooGuFilter.InsertUnique([]byte(short)); err != false {
	// 	logx.Errorw("l.svcCtx.BooGuFilter.InsertUnique", logx.LogField{Key: "err", Value: "not exit"})
	// 	//生产建议：如果布隆过滤器报错（比如 Redis 挂了），为了业务可用性，通常选择“放行”而不是报错
	// }

	// 5. 返回响应
	// 5.1 返回的是短域名+短链接
	shortUrl := l.svcCtx.Config.ShortDomain + "/" + short
	return &types.ConvertResponse{ShortURL: shortUrl}, nil
}
