package servers

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/metux/go-nebulon/filestore"
)

//go:embed assets/video.html
var asset_video_html string

func videoInitAssets(router *gin.Engine) {
	t := template.New("video.html")
	t.Parse(asset_video_html)
	router.SetHTMLTemplate(t)
}

func (server *HttpServer) serveGetVideo(ctx *gin.Context) {
	fs := filestore.NewFileStore(server.Conf.BlockStore)

	ref, ok := server.getRefParam(ctx)
	if !ok {
		return
	}

	fctrl, err := fs.ReadFileControl(ref)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	contenttype, _ := fctrl.Header["Content-Type"]

	video_url := fmt.Sprintf(FileURLPrefix+"%s/%s/%s/%s",
		ctx.Param("reftype"),
		ctx.Param("id"),
		ctx.Param("cipher"),
		ctx.Param("key"))

	ctx.HTML(http.StatusOK, "video.html", gin.H{
		"video_url":      video_url,
		"video_mimetype": contenttype,
		"video_id":       ctx.Param("id")})
}
