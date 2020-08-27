package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paramSearch struct {
	Offset    int    `form:"offset"`
	Length    int    `form:"length"`
	SortKey   string `form:"sortKey"`
	SortOrder string `form:"sortOrder"`
}

type paramId struct {
	Id int64 `uri:"id"`
}

type paramId2 struct {
	Id  int64 `uri:"id"`
	Id2 int64 `uri:"id2"`
}

func RegisterRoutes(app *gin.RouterGroup) {

	app.GET("/channels", channelList)
	app.POST("/channels", channelList)
	app.POST("/channel", channelCreate)
	app.DELETE("/channel/:id", channelDelete)
	app.PUT("/channel/:id", channelEdit)
	app.GET("/channel/:id", channelGet)

}

func responseOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": data,
	})
}

func responseError(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":    false,
		"error": err,
	})
}

func nop(c *gin.Context) {
	c.String(http.StatusForbidden, "Unsupported")
}

func mustLogin(c *gin.Context) {
	//测试
	session := sessions.Default(c)
	if user := session.Get("user"); user != nil {
		c.Next()
	} else {
		//c.Redirect(http.StatusSeeOther, "/login")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
