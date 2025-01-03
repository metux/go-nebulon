package httpd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) MP4File(ctx *gin.Context) {
	ranges, err := parseRange(ctx.Request.Header.Get("Range"), 1024*1024*1024*4)

	log.Printf("ranges: %+v\n", ranges)
	log.Printf("ranges err: %s\n", err)

	reader, headers, err := server.fs.ReadStream(server.Ref)

	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	if reader == nil {
		ctx.JSON(http.StatusNotFound, "no reader")
	}

	//	log.Printf("Request headers: %+v\n", ctx.Request.Header)

	//	rangeheader := ctx.Request.Header.Get("Range")
	//	log.Printf("range header: %s\n", rangeheader)

	ctx.DataFromReader(
		http.StatusOK,
		-1,
		headers["Content-Type"],
		reader,
		headers)
}
