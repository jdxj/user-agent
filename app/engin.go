package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jdxj/user-agent/control"

	"github.com/astaxie/beego/logs"

	"github.com/gin-gonic/gin"
	"github.com/jdxj/user-agent/db"
	"github.com/jdxj/user-agent/module"
)

func (coll *Collector) newEngine() *gin.Engine {
	r := gin.Default()
	r.Use(coll.RejectFaviconIco)
	r.Use(coll.RecordHeader)
	r.Any("/", Ping)

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
	coll.headerInfos <- headerInfo
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

func (coll *Collector) cacheHeaderInfo() {
	headerInfos := make([]*module.HeaderInfo, 0, headerInfoCacheLimit)

	for {
		select {
		case <-control.Stop:
			headerInfos = headerInfos[:0]
			logs.Debug("stop cache headerInfo")
			return

		case hi := <-coll.headerInfos:
			headerInfos = append(headerInfos, hi)
		}

		if len(headerInfos) >= headerInfoCacheLimit {
			InsertHeader(headerInfos)
			headerInfos = headerInfos[:0]
		}
	}
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (coll *Collector) RejectFaviconIco(c *gin.Context) {
	path := c.Request.RequestURI
	if strings.Index(path, "favicon.ico") >= 0 {
		c.AbortWithStatus(http.StatusNotFound)
		logs.Debug("abort favicon, request url: %s", path)
	}
}
