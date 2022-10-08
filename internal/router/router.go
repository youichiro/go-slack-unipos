package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.IndentedJSON(200, gin.H{"message": "hello world"}) })

	return r
}
