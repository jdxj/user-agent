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

	// middleware
	r.Use(coll.RejectFaviconIco)
	r.Use(coll.RejectEmptyUserAgent)
	r.Use(coll.RecordHeader)

	// handler
	r.Any("/", Ping)
	r.Any("/index.html", Ping)
	r.Any("/smb_scheduler/cdr.htm", Ping)
	r.Any("/goip/cron.htm", Ping)
	r.Any("/navigation.html", Ping)
	r.Any("/KingViewWeb", Ping)
	r.Any("/webconfig.ini", Ping)
	r.Any("/echo.php", Ping)
	r.Any("/cgi-bin/mainfunction.cgi", Ping)
	r.Any("/xsser.php", Ping)
	r.Any("/forums/index.php", Ping)
	r.Any("/bbs/index.php", Ping)
	r.Any("/license.php", Ping)
	r.Any("/v/index.php", Ping)
	r.Any("/s/index.php", Ping)
	r.Any("/1/index.php", Ping)
	r.Any("/adv,/cgi-bin/weblogin.cgi", Ping)
	r.Any("/nice%20ports%2C/Tri%6Eity.txt%2ebak", Ping)
	r.Any("/sdk", Ping)
	r.Any("/nmaplowercheck1587934706", Ping)
	r.Any("/evox/about", Ping)
	r.Any("/HNAP1", Ping)
	r.Any("/index.php", Ping)
	r.Any("/manager/text/list", Ping)
	r.Any("/sqlite/main.php", Ping)
	r.Any("/sqlitemanager/main.php", Ping)
	r.Any("/SQLiteManager/main.php", Ping)
	r.Any("/SQLite/main.php", Ping)
	r.Any("/main.php", Ping)
	r.Any("/test/sqlite/SQLiteManager-1.2.0/SQLiteManager-1.2.0/main.php", Ping)
	r.Any("/SQLiteManager-1.2.4/main.php", Ping)
	r.Any("/agSearch/SQlite/main.php", Ping)

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

func (coll *Collector) RejectEmptyUserAgent(c *gin.Context) {
	if c.GetHeader("User-Agent") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "must user-agent",
		})
	}
}
