package gocache

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type RestServer struct {
	cache Cache
	r     *gin.Engine
}

func NewRestServer() *RestServer {
	gin.SetMode(gin.ReleaseMode)
	server := &RestServer{cache: newMemoryCache(), r: gin.Default()}

	server.r.PUT("/cache/:key", server.handleSetCache)
	server.r.DELETE("/cache/:key", server.handleDelCache)
	server.r.GET("/cache/:key", server.handleGetCache)
	server.r.GET("/status", server.handleCacheStat)
	return server
}

func (server *RestServer) Run(addr string) {
	server.r.Run(addr)
}

func (server *RestServer) handleSetCache(c *gin.Context) {
	value, _ := ioutil.ReadAll(c.Request.Body)
	if len(value) == 0 {
		c.Status(500)
		return
	}

	server.cache.Set(c.Param("key"), value)
	c.Status(200)
}

func (server *RestServer) handleGetCache(c *gin.Context) {
	value, err := server.cache.Get(c.Param("key"))
	if err != nil {
		c.Status(500)
		return
	}
	c.Data(200, "plain", value)
}

func (server *RestServer) handleCacheStat(c *gin.Context) {
	c.AsciiJSON(200, server.cache.GetStat())
}

func (server *RestServer) handleDelCache(c *gin.Context) {
	server.cache.Del(c.Param("key"))
}
