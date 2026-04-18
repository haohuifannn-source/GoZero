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
	// 2. 取号
	// 3. 号码转短链
	// 4. 存储长链接和短链接的映射关系
	// 5. 返回响应
	return
}
