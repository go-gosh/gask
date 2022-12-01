package resource

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static
var file embed.FS

//go:embed index.html
var index []byte

//go:embed favicon.ico
var favicon []byte

type r struct {
}

func (r r) Open(name string) (fs.File, error) {
	return file.Open("static/" + name)
}

func Setup(engine *gin.Engine) {
	engine.StaticFS("/static", http.FS(r{}))
	engine.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	engine.Any("/favicon.ico", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/plain", favicon)
	})
}
