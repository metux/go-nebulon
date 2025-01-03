package httpd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/wire"
)

type Server struct {
	Router *gin.Engine
	bs     base.BlockStore
	fs     base.FileStore
	Ref    wire.BlockRef
}

func NewServer(bs base.BlockStore, fs base.FileStore) *Server {
	s := new(Server)
	s.bs = bs
	s.fs = fs
	s.Router = gin.Default()
	s.Router.GET("/", s.HomePage)
	s.Router.GET("/mp4", s.MP4File)
	s.Router.GET("/block/:id", s.GetBlock)
	log.Printf("router initialized\n")
	return s
}

func (s *Server) DoUpload(fn string, ct string) {
	log.Printf("uploading file %s\n", fn)
	ref, err := helpers.PutFile(s.fs,
		"",
		wire.Header{wire.Header_ContentType: ct},
		fn)
	if err != nil {
		panic(err)
	}
	s.Ref = ref
	log.Printf("uploaded file: %s\n", ref.Dump())
	log.Printf("uploaded file: %s\n", s.Ref.Dump())
}

func (s *Server) Run(addr string) {
	log.Printf("starting up web server\n")
	s.Router.Run(addr)
}
