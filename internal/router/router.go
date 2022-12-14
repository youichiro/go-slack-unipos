package router

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/youichiro/go-slack-my-unipos/internal/handler"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	slackHandler := handler.SlackHandler{
		Db:           db,
		SigninSecret: os.Getenv("SLACK_SIGNING_SECRET"),
		Token:        os.Getenv("SLACK_BOT_USER_OAUTH_TOKEN"),
	}

	r.GET("/", func(c *gin.Context) { c.IndentedJSON(200, gin.H{"message": "hello world"}) })
	r.POST("/slash", slackHandler.HandleSlash)
	r.POST("/modal", slackHandler.HandleModal)

	return r
}
