// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	//"github.com/redis/go-redis/v9"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"github.com/zeromicro/go-zero/core/bloom"
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

	// 布隆过滤器
	FilterBloom *bloom.Filter

	// 布谷鸟过滤器
	BooGuFilter *cuckoo.Filter
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

	// 一般来说都要不一样，这里为了方便，直接用CacheRedis的地址
	stroe := redis.New(c.CaCheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})

	// 声明一个bitSet
	filter := bloom.New(stroe, "bloom_filter", 20*1<<20)

	// 初始化一个布谷鸟过滤器
	cf := cuckoo.NewFilter(1000)

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(shortConn, c.CaCheRedis),
		SquenceModel:      sequence.NewMySql(c.Sequence.DSN, "a"),
		SquenceRedis:      sequence.NewRedis(c.SequenceRedis.Host, "sequence:a:"),
		ShortUrlBlackList: m,
		BizRedis:          bizRedis,
		FilterBloom:       filter,
		BooGuFilter:       cf,
	}

}

// loadDataToBloomFilter记载已有的短链接数据至布隆过滤器
// 要实现内存版本的需要实现这个函数，在每次程序运行的时候都需要先去加载
// func loadDataToBloomFilter(c config.Config, filter *bloom.Filter) {
//     // 1. 初始化数据库连接（使用 go-zero 的 sqlx）
//     shortConn := sqlx.NewMysql(c.ShortUrlDB.DSN)

//     // 2. 这里我们不一定要用生成的 Model，直接用原始 SQL 查出所有 surl 效率最高
//     // 我们只需要 surl 这一列，不需要整行数据
//     query := "SELECT surl FROM short_url_map"

//     // 3. 定义一个变量接收数据
//     var surls []string

//     // 4. 执行查询
//     // 注意：如果数据量巨大（千万级），建议使用分页查询或 stream 处理
//     err := shortConn.QueryRowsNoCache(&surls, query)
//     if err != nil {
//         logx.Errorf("Load data to bloom filter failed: %v", err)
//         return
//     }

//     // 5. 将数据逐个添加到布隆过滤器中
//     count := 0
//     for _, surl := range surls {
//         if surl != "" {
//             _ = filter.Add([]byte(surl))
//             count++
//         }
//     }

//     logx.Infof("Bloom filter pre-warming finished. Loaded %d keys.", count)
// }
