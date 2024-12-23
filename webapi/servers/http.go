package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/webapi/blockstore/rest_v1"
)

const (
	FileURLPrefix   = "/file/"
	PlayerURLPrefix = "/video/"
)

type HttpServer struct {
	Conf   base.ServerConfig
	Router *gin.Engine
}

func (server *HttpServer) Shutdown(force bool) error {
	if force {
		//                server.RPCServer.Stop()
	} else {
		//                server.RPCServer.GracefulStop()
	}
	return nil
}

func (server *HttpServer) homePage(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to the Nebulon REST API")
}

func (server *HttpServer) Serve() error {
	return server.Router.Run(":" + server.Conf.Port)
}

func NewHTTPServer(cf base.ServerConfig) *HttpServer {
	server := HttpServer{
		Conf:   cf,
		Router: gin.Default(),
	}
	videoInitAssets(server.Router)
	server.Router.GET("/", server.homePage)
	server.Router.GET(FileURLPrefix+":reftype/:id/:cipher/:key", server.serveGetFile)
	server.Router.GET(PlayerURLPrefix+":reftype/:id/:cipher/:key", server.serveGetVideo)
	v1 := server.Router.Group("/v1")
	rest_v1.Register(v1, server.Conf.BlockStore)

	return &server
}
