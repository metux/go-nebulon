package httpd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) MP4File(ctx *gin.Context) {
	reader, headers, err := server.fs.ReadStream(server.Ref)

	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	if reader == nil {
		ctx.JSON(http.StatusNotFound, "no reader")
	}

	log.Printf("Request headers: %+v\n", ctx.Request.Header)

	rangeheader := ctx.Request.Header.Get("Content-Range")
	log.Printf("range header: %s\n", rangeheader)

	ctx.DataFromReader(
		http.StatusOK,
		-1,
		headers["Content-Type"],
		reader,
		headers)
}
