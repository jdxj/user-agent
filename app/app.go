package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jdxj/user-agent/control"
	"github.com/jdxj/user-agent/db"

	"github.com/astaxie/beego/logs"

	"github.com/jdxj/user-agent/module"
)

const headerInfoCacheLimit = 5
const address = ":80"

func NewCollector() *Collector {
	coll := &Collector{
		headerInfos: make(chan *module.HeaderInfo, headerInfoCacheLimit),
	}

	srv := &http.Server{
		Addr:    address,
		Handler: coll.newEngine(),
	}
	coll.srv = srv

	return coll
}

type Collector struct {
	srv *http.Server

	// todo: 锁
	headerInfos chan *module.HeaderInfo
}

func (coll *Collector) Start() {
	srv := coll.srv

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Error("web listen err: %s", err)
			close(control.Stop)
			db.MySQL.Close()
			logs.GetBeeLogger().Flush()
			panic(err)
		}
	}()

	go coll.cacheHeaderInfo()
}

func (coll *Collector) Stop() {
	srv := coll.srv
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logs.Error("web shutdown err: %s", err)
		return
	}

	logs.Debug("wait for stop web")
	<-ctx.Done()
}
