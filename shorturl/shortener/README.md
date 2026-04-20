# 短链接项目

## 搭建项目的骨架

1. 建库建表

新建发射器表
```sql
CREATE TABLE `sequence` (
 `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
 `stub` varchar(1) NOT NULL,
 `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 PRIMARY KEY (`id`),
 UNIQUE KEY `idx_uniq_stub` (`stub`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT = '序号表';
```

新建短链接和长链接的映射表，这里通过长链接的MD5码做了唯一索引
```sql
CREATE TABLE `short_url_map` (
 `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
 `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
 `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
 `is_del` tinyint UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',
 
 `lurl` varchar(2048) DEFAULT NULL COMMENT '⻓长链接',
 `md5` char(32) DEFAULT NULL COMMENT '⻓长链接MD5', 
 `surl` varchar(11) DEFAULT NULL COMMENT '短链接',
  PRIMARY KEY (`id`),
  INDEX(`is_del`),
  UNIQUE(`md5`),
  UNIQUE(`surl`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '⻓长短链映射表';
```

2. 搭建go-zero框架的骨架


编写`.api`文件，使用goctl命令生成代码
```api
goctl api go -api shortener.api -dir . -style=goZero
```

3. 根据数据表生成model层代码
```api
goctl model mysql datasource -url="root:root@tcp(117.72.109.40:3306)/db1" --table="short_url_map" --dir="./model" --style="goZero"

goctl model mysql datasource -url="root:root@tcp(117.72.109.40:3306)/db1" --table="sequence" --dir="./model" --style="goZero"
```

4. 下载项目依赖
```bash
go mod tidy
```

5. 运行启动项目
```bash
go run shortner.go 
```

6. 设置配置文件
注意：两边一定要对齐！！！！

## 实现长链接转短链接的功能（conver）

1. 校验输入的数据

1.1 数据不可以为空 / 长链接能请求通的网址 / 判断数据库中是否已经存在长链接 / 输入的不能是一个短链接（防止循环转链）

(1) 数据不可以为空：采用go知名的第三方库validator实现，防止过多的if else分支:https://golang.halfiisland.com/community/pkgs/validate/Validator.html#%E4%BB%8B%E7%BB%8D

```bash
go get github.com/go-playground/validator/v10

import "github.com/go-playground/validator/v10"
```

(2) 长链接能请求通的网址/判断数据库中是否已经存在长链接/输入的不能是一个短链接（防止循环转链）:都通过封装为pkg来实现

***注意***: (1) 判断数据库中是否已经存在长链接中是通过长链接转化为MD5码进行的，其中hex.EncodeToString(h.Sum(nil))参数nil 表示不需要将结果追加到现有的切片中，直接返回一个新的 []byte；(2) 在输入的不能是一个短链接中，basePath := path.Base(myUrl.Path)是可以得到其中的path，https://www.example.com/posts/123?auth=true中Scheme (https), Host (www.example.com), Path (/posts/123) 等部分。因此可以解析后直接得到path部分

1.1.2 为结构体添加validata tag， 并添加校验规则

(1) 在handler函数中校验

2. 取号

***重点***: 这里面是通过创建一个sequence进行next函数的编写，这里面定义了一个sequence的空接口类型，是为了无论是mysql实现还是redis实现，只要实现了接口方法，就可以调用，不需要再svc里面多次初始化，只需要将其定义为sequence.Squence类型即可。

这里实现了两种方法，mysql和reids。

3. 号码转短链

***注意***: 这里面通过创建base62实现了62进制的取号操作，这里的base62Str通过配置文件去指定，可以防止别人恶意请求去扒数据库，同时在主函数里面去初始化这个字符串。同时加入了黑名单等，防止转链后的路径具有非法的名称。黑名单的操作通过转换map的存储方法使得避免for循环的引入，同时采用map[stirng]struct{}的方式是因为空结构体不占用内存。

4. 存储长链接和短链接的映射关系

5. 返回响应

3. 编写单元测试
(1) 方法一：可以通过编译器自动生成单元测试，然后修改相关的代码，验证自己所实现的功能的正确性，该项目针对了urltool.go的代码正确性。

（2）方法二：通过goconvey方法https://liwenzhou.com/posts/go/unit-test-5/
```bash
go install github.com/smartystreets/goconvey@latest
```

然后编写即可，如connect_test.go所示，然后可以通过命令运行所有的测试
```bash
go test ./...
```

## 实现查看短链的功能

(1) 通过数据库查询短链对应的长链

(2) 为了实现重定向，需要去handler层实现重定向

(3) 为了实现功能的高可用性，查数据库之前必须加缓存，防止大规模的查询数据库，有两种方式
1. 自己写缓存，就是可以是按surl -->  lurl
2. go-zero生成的缓存，就是 surl ---> 数据行，即将模型的结构的数据全部缓存， 会增加缓存的压力

为了方便，这里用第二种方式：
1. 添加缓存的配置.yaml和config结构体都需要改
2. 删除旧的model层代码，删除shorturlmapmodel.go文件，需要改哪个就删那个
3. 重新生成model层代码
```bash
goctl model mysql datasource -url="root:root@tcp(117.72.109.40:3306)/db1" --table="short_url_map" --dir="./model" --style="goZero" -c
```

第一种也实现了：
1. 首先修改配置文件和conf文件，加入redis的相关配置
2. 修改svc文件，加入redis
3. 在logic除调用，主要逻辑是首先去查redis，查到了返回结果；查不到就去查数据库，然后同步设置redis

4. 修改svc代码

(4) 使用Redis会存在几个问题
1. 缓存如何设置，LRU
2. 缓存击穿：缓存过期了，大量的请求同一时间穿透去查数据库
- 使⽤singleflight 合并请求：https://www.liwenzhou.com/posts/Go/singleflight/
- go-zero本身就支持singleflight，带有缓存的model生成的时候就带有这个

```go
//svc文件
// internal/svc/servicecontext.go
import (
    "github.com/zeromicro/go-zero/core/syncx" // 引入 syncx
)

type ServiceContext struct {
    // ... 其他字段
    SingleGroup syncx.SingleFlight // 添加这个字段
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        // ...
        SingleGroup: syncx.NewSingleFlight(), // 初始化
    }
}

// 防止缓存击穿的手动实现版
// 缓存击穿一定是发生在缓存查不到的情况下，因为一开始肯定是先要去查，再做缓存击穿防御
// 如果不想用dochan版本，只需要把select和dochan改为do即可
func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
    // 1. 构造 Redis Key
    redisKey := fmt.Sprintf("%s%s", reidisDb1ShortUrlMapSurlPrefix, req.ShortURL)

    // 2. 第一层：尝试从手动缓存读取
    s, _ := l.svcCtx.BizRedis.Get(redisKey)
    if len(s) > 0 {
        return &types.ShowResponse{LongURL: s}, nil
    }

    // 3. 第二层：使用 DoChan 进行异步归并
    // 注意：DoChan 不会阻塞，它会立刻返回一个 channel
    resultChan := l.svcCtx.SingleGroup.DoChan(req.ShortURL, func() (any, error) {
        // --- 依然只有一个人会进来 ---
        u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(
            l.ctx, 
            sql.NullString{String: req.ShortURL, Valid: true},
        )
        if err != nil {
            return nil, err
        }

        // 回填缓存
        _ = l.svcCtx.BizRedis.Setex(redisKey, u.Lurl.String, 86400)
        return u.Lurl.String, nil
    })

    // 4. 使用 select 监听结果或超时
    select {
    case <-l.ctx.Done():
        // 情况 A：如果在数据库返回前，客户端断开了连接或请求超时了
        return nil, l.ctx.Err()

    case res := <-resultChan:
        // 情况 B：数据库返回了结果
        if res.Err != nil {
            if res.Err == sqlx.ErrNotFound {
                return nil, errors.New("404")
            }
            return nil, res.Err
        }
        
        // res.Val 是 any 类型，断言为 string
        return &types.ShowResponse{LongURL: res.Val.(string)}, nil
    }
}
```

3. 缓存穿透：根本就没有这个链接的缓存，而且数据库也没有，因此也会大量请求并发的查询数据库


***注意***： 在convert中通过拼接得到q1mi.cn/L，其实q1mi.cn是个人购买的域名，在本地开发中就相当于127.0.0.1:8888，因为是只需要在网址输入localhost:8888/L就可以去实现功能的校验