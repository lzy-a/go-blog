package v1

import (
	"gin/models"
	"gin/pkg/e"
	"gin/pkg/setting"
	"gin/pkg/util"
	"log"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// //获取文章列表
// apiv1.GET("/articles", v1.GetArticles)
// //获取指定文章
// apiv1.GET("/articles/:id", v1.GetArticle)
// //新增文章
// apiv1.POST("/articles", v1.AddArticle)
// //修改文章
// apiv1.PUT("/articles/:id", v1.EditArticle)
// //删除文章
// apiv1.DELETE("/articles/:id", v1.DeleteArticle)

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			data = models.GetArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// 查询某tagID对应的文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	var tagID int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagID, _ = strconv.Atoi(arg)
		maps["tag_id"] = tagID
		valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	}
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticlesTotal(maps)
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Query("tag_id"))
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state, _ := strconv.Atoi(c.DefaultQuery("state", "0"))

	valid := validation.Validation{}
	valid.Min(tagID, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagID) {
			data := make(map[string]interface{})
			data["tag_id"] = tagID
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	tagID, _ := strconv.Atoi(c.Query("tag_id"))
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)

	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Required(tagID, "tag_id").Message("标签ID不能为空")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			data := make(map[string]interface{})
			if tagID > 0 {
				data["tag_id"] = tagID
			}
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["modified_by"] = modifiedBy
			data["state"] = state
			models.EditArticle(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
		}
	}
	c.JSON(code, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
