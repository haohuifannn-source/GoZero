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

3. 号码转短链

4. 存储长链接和短链接的映射关系

5. 返回响应
