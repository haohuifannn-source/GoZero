package sequence

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 非硬编码
const sqlxReplaceIntoStub = `REPLACE INTO sequence (stub) VALUES (?);`

type Mysql struct {
	conn sqlx.SqlConn
	sub  string
}

func NewMySql(connStr string, subStr string) Sequence {
	return &Mysql{
		conn: sqlx.NewMysql(connStr),
		sub:  subStr,
	}
}

func (m *Mysql) Next() (seq uint64, err error) {
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlxReplaceIntoStub)
	if err != nil {
		logx.Errorw("m.conn.Prepare failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	defer stmt.Close()

	rest, err := stmt.Exec(m.sub)
	if err != nil {
		logx.Errorw("stmt.Exec() failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}

	lid, err := rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId() failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}

	return uint64(lid), nil

}
