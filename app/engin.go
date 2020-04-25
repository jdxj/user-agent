package app

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"

	"github.com/gin-gonic/gin"
	"github.com/jdxj/user-agent/db"
	"github.com/jdxj/user-agent/module"
)

func (coll *Collector) newEngine() *gin.Engine {
	r := gin.Default()
	r.Any("/", coll.RecordHeader)

	return r
}

func (coll *Collector) RecordHeader(c *gin.Context) {
	req := c.Request

	headerInfo := &module.HeaderInfo{
		IP:        req.RemoteAddr,
		Host:      req.Host,
		Referer:   req.Referer(),
		UserAgent: req.UserAgent(),
		Method:    req.Method,
		Path:      req.RequestURI,
	}
	coll.headerInfos = append(coll.headerInfos, headerInfo)

	if len(coll.headerInfos) >= headerInfoCacheLimit {
		InsertHeader(coll.headerInfos)
		coll.headerInfos = coll.headerInfos[:0]
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"len":     len(coll.headerInfos),
	})
}

func InsertHeader(headerInfos []*module.HeaderInfo) {
	mysql := db.MySQL

	query := fmt.Sprintf("INSERT INTO request (ip,host,referer,user_agent,method,path) VALUES (?,?,?,?,?,?)")
	stmt, err := mysql.Prepare(query)
	if err != nil {
		logs.Error("stmt err: %s", err)
		return
	}
	defer stmt.Close()

	for _, hi := range headerInfos {
		_, err := stmt.Exec(hi.IP, hi.Host, hi.Referer, hi.UserAgent, hi.Method, hi.Path)
		if err != nil {
			logs.Error("stmt exec err: %s", err)
			logs.Error("data: %#v", *hi)
			continue
		}
	}
}
