package servers

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
	"github.com/metux/go-nebulon/filestore"
	"github.com/metux/go-nebulon/util"
)

func (server HttpServer) getRefParam(ctx *gin.Context) (base.BlockRef, bool) {
	p_oid := ctx.Param("id")
	p_key := ctx.Param("key")

	oid, err := hex.DecodeString(p_oid)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, base.ErrInvalidOID)
		return base.BlockRef{}, false
	}

	key, err := hex.DecodeString(p_key)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, base.ErrInvalidKey)
		return base.BlockRef{}, false
	}

	return base.BlockRef{
		Oid:    oid,
		Type:   wire.ParseRefType(ctx.Param("reftype")),
		Cipher: wire.ParseCipherType(ctx.Param("cipher")),
		Key:    key}, true
}

// FIXME: support limit (range end)
func (server *HttpServer) serveGetFile(ctx *gin.Context) {
	fs := filestore.NewFileStore(server.Conf.BlockStore)

	ref, ok := server.getRefParam(ctx)
	if !ok {
		return
	}

	replyStatus := http.StatusOK

	ranges, err := util.ParseRangeHeader(ctx.Request.Header)
	if err != nil {
		ctx.String(http.StatusBadRequest, "broken range header: %s\n", ctx.Request.Header["Range"])
		return
	}

	if len(ranges) > 1 {
		ctx.String(http.StatusBadRequest, "multiple ranges not supported: %s\n", ctx.Request.Header["Range"])
		return
	}

	reader, headers, size, err := fs.ReadStream(ref, uint64(ranges[0].StartPos))
	defer reader.Close()

	ctx.Header("Content-Length", fmt.Sprintf("%d", size))

	if len(ranges) == 1 {
		replyStatus = http.StatusPartialContent
		ctx.Header("Accept-Ranges", "bytes")
		ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ranges[0].StartPos, size-1, size))
	}

	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	if reader == nil {
		ctx.JSON(http.StatusNotFound, "no reader")
	}

	ctx.DataFromReader(
		replyStatus,
		-1,
		headers["Content-Type"],
		reader,
		headers)
}
