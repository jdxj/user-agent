package app

import (
	"testing"

	"github.com/jdxj/user-agent/control"
	"github.com/jdxj/user-agent/db"

	"github.com/jdxj/user-agent/module"
)

func TestInsertHeader(t *testing.T) {
	hi := &module.HeaderInfo{
		IP:        "ip",
		Host:      "host",
		Referer:   "referer",
		UserAgent: "userAgent",
	}
	InsertHeader([]*module.HeaderInfo{hi})

	close(control.Stop)
	db.MySQL.Close()
}
