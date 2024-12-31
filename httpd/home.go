package httpd

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Server) HomePage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the Task Manager API"})
}
