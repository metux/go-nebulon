package httpd

import (
	"log"

	"github.com/gin-gonic/gin"
	//	"github.coOBm/metux/xsx-middleware/dbxs"
	//	"github.com/metux/xsx-middleware/secdocs"
	"github.com/metux/go-nebulon/base"
	"github.com/metux/go-nebulon/helpers"
	"github.com/metux/go-nebulon/wire"
)

type Server struct {
	//	dbx      dbxs.DBXS
	Router *gin.Engine
	fs     base.FileStore
	//	sdclient secdocs.Client
	Ref wire.BlockRef
}

func NewServer(fs base.FileStore) *Server {
	s := new(Server)
	s.fs = fs
	s.Router = gin.Default()
	s.Router.GET("/", s.HomePage)
	s.Router.GET("/mp4", s.MP4File)
	//	s.Router.GET("/users/:user/boxes", s.ListBoxes)
	//	s.Router.PUT("/upload", s.Upload)
	//	s.Router.GET("/aoid-evidence/:aoid", s.AOIDEvidence)
	//	s.Router.GET("/aoid-retrieve/:aoid", s.AOIDRetrieve)
	log.Printf("router initialized\n")
	return s
}

func (s *Server) DoUpload(fn string, ct string) {
	log.Printf("uploading file %s\n", fn)
	ref, err := helpers.StoreFile(s.fs,
		map[string]string{"Content-Type": ct},
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
