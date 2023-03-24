package routers

import (
	"gin/pkg/setting"
	v1 "gin/routers/api/v1"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	router := gin.Default()
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "test",
		})
	})
	gin.SetMode(setting.RunMode)

	apiv1 := router.Group("/api/v1")
	{
		//获取文章标签列表
		apiv1.GET("/tags", v1.GetTags)
		//添加标签
		apiv1.POST("/tags", v1.AddTags)
		//修改标签
		apiv1.PUT("/tags/:id", v1.EditTags)
		//删除标签
		apiv1.DELETE("/tags/:id", v1.DeleteTags)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新增文章
		apiv1.POST("/articles", v1.AddArticle)
		//修改文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return router
}