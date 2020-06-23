package db

import (
	"database/sql"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/jdxj/user-agent/control"

	_ "github.com/go-sql-driver/mysql"
)

const pingDur = 5 * time.Minute

var MySQL = newMySQL()

func newMySQL() *sql.DB {
	dsn := ":@tcp(127.0.0.1:49160)/http?loc=Local&parseTime=true"
	//dsn := ":@@tcp(mysql.aaronkir.xyz:49160)/http?loc=Local&parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(pingDur)
		defer ticker.Stop()

		for {
			select {
			case <-control.Stop:
				logs.Debug("stop db ping")
				return

			case <-ticker.C:
			}

			if err := db.Ping(); err != nil {
				logs.Error("db ping err: %s", err)
				return
			}
			logs.Debug("db ping success")
		}
	}()

	return db
}
