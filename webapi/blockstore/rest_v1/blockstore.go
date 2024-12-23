package rest_v1

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

type HttpBlockStore struct {
	BlockStore base.IBlockStore
}

func (server HttpBlockStore) replyRef(ctx *gin.Context, ref base.BlockRef) {
	ctx.Header(wire.HttpV1_Header_BlockRef, ref.Dump())
}

func (server HttpBlockStore) getRef(ctx *gin.Context) (base.BlockRef, bool) {
	p_type := ctx.Param("reftype")
	p_oid := ctx.Param("id")

	oid, err := hex.DecodeString(p_oid)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid oid [%w]", err))
		return base.BlockRef{}, false
	}

	if val, ok := wire.RefType_value[p_type]; ok {
		return base.BlockRef{Oid: oid, Type: wire.RefType(val)}, true
	}

	ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid reftype [%s]", p_type))
	return base.BlockRef{Oid: oid}, false
}

// swagger:route GET /block/{RefType}/{RefOID} blocks ServeBlockGet
//
// retrieve a specific block (by ref type and ref oid)
// responds with raw block data as application/octet-stream
//
// Responses:
//
//	200: blockDataResponse
//	400: badRequest
//	404: notFoundResponse
//	500: internalErrorResponse
func (server HttpBlockStore) ServeBlockGet(ctx *gin.Context) {
	ref, ok := server.getRef(ctx)
	if !ok {
		return
	}
	server.replyRef(ctx, ref)

	data, err := server.BlockStore.GetBlock(ref)

	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("failed loading block [%w]", err))
		return
	}

	ctx.DataFromReader(
		http.StatusOK,
		-1,
		"application/octet-stream",
		bytes.NewReader(data),
		map[string]string{})
}

// swagger:route PUT /block/{RefType}/{RefOID} blocks ServeBlockPut
//
// upload binary data block (by ref type and ref oid)
// responds with final object ref
//
// Responses:
//
//	200: objectCreatedResponse
//	400: badRequest
//	404: notFoundResponse
//	406: notAcceptableResponse
//	500: internalErrorResponse
func (server HttpBlockStore) ServeBlockPut(ctx *gin.Context) {
	ref, ok := server.getRef(ctx)
	if !ok {
		return
	}
	server.replyRef(ctx, ref)

	raw, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("invalid raw data [%w]", err))
		return
	}

	newref, err := server.BlockStore.PutBlock(raw, ref.Type)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed storing block [%w]", err))
		return
	}

	if bytes.Equal(ref.Oid, newref.Oid) {
		ctx.String(http.StatusCreated, newref.Dump())
		return
	}

	ctx.AbortWithError(http.StatusNotAcceptable, fmt.Errorf("oids mismatch: want %v stored %v\n", ref.Oid, newref.Oid))
}

// swagger:route PATCH /block/{RefType}/{RefOID} blocks ServeBlockKeep
//
// mark block to be kept longer
//
// Responses:
//
//	201: keepingBlock
//	400: badRequest
//	404: notFoundResponse
//	500: internalErrorResponse
func (server HttpBlockStore) ServeBlockKeep(ctx *gin.Context) {
	ref, ok := server.getRef(ctx)
	if !ok {
		return
	}
	server.replyRef(ctx, ref)

	if err := server.BlockStore.KeepBlock(ref); err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("failed touching block [%w]", err))
		return
	}

	ctx.String(http.StatusNoContent, "keeping: "+ref.Dump())
}

func (server HttpBlockStore) ServeBlockInfGet(ctx *gin.Context) {
	ref, ok := server.getRef(ctx)
	if !ok {
		return
	}
	server.replyRef(ctx, ref)

	depth, _ := strconv.Atoi(ctx.Request.Header.Get(wire.HttpV1_Header_FetchDepth))
	if inf, err := server.BlockStore.PeekBlock(ref, depth); err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("failed peeking block [%w]", err))
		return
	} else {
		ctx.JSON(http.StatusOK, inf)
	}
}

// swagger:route GET /block blocks ServeBlockList
//
// list all stored object references
//
// Responses:
//
//	201: blockrefListResponse
func (server HttpBlockStore) ServeBlockList(ctx *gin.Context) {
	ch := server.BlockStore.IterateBlocks()
	ctx.Stream(func(w io.Writer) bool {
		if ent, ok := <-ch; ok {
			if !ent.Finished && ent.Error == nil {
				w.Write([]byte(ent.Ref.Dump() + "\n"))
			}
			return true
		}
		return false
	})
}

// swagger:route GET / main ServeHello
//
// say hello
//
// Responses:
//
//	201: helloResponse
func (server HttpBlockStore) ServeHello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "nebulon blockstore http api v1")
}

func (server HttpBlockStore) AddRoutes(router gin.IRoutes) {
	// FIXME: add HEAD request for Peek
	router.GET("/", server.ServeHello)
	router.GET("/block", server.ServeBlockList)
	router.GET("/block/:reftype/:id", server.ServeBlockGet)
	router.PUT("/block/:reftype/:id", server.ServeBlockPut)
	router.PATCH("/block/:reftype/:id", server.ServeBlockKeep)
	router.GET("/blockinf/:reftype/:id", server.ServeBlockInfGet)
}

func Register(router gin.IRoutes, bs base.IBlockStore) {
	server := HttpBlockStore{BlockStore: bs}
	server.AddRoutes(router)
}
