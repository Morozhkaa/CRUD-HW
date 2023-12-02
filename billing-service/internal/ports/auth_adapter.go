package ports

import "github.com/gin-gonic/gin"

type AuthAdapter interface {
	Verify(ctx *gin.Context) error
}
