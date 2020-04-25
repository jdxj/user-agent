package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/user-agent/app"
	"github.com/jdxj/user-agent/control"
	"github.com/jdxj/user-agent/db"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"user-agent.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	coll := app.NewCollector()
	coll.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		logs.Debug("receive stop signal")
	}

	coll.Stop()
	close(control.Stop)
	db.MySQL.Close()

	logs.Debug("stop app")
	logs.GetBeeLogger().Flush()
}
