package httpd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metux/go-nebulon/wire"
)

func (server *Server) GetBlock(ctx *gin.Context) {
	oid, err := hex.DecodeString(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, fmt.Errorf("invalid oid [%w]", err))
		return
	}

	ref := wire.BlockRef{Oid: oid}
	data, err := server.bs.LoadBlock(ref)

	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	ctx.DataFromReader(
		http.StatusOK,
		-1,
		"", /* headers["Content-Type"], */
		bytes.NewReader(data),
		map[string]string{})
}
