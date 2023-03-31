package v1

import (
	"gin/pkg/app"
	"gin/pkg/e"
	"gin/pkg/export"
	"gin/pkg/setting"
	"gin/pkg/util"
	tagservice "gin/service/tag_service"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} string "ok"
// @Failure 500 {object} string "fail"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{c}
	name := c.Query("name")
	data := make(map[string]interface{})
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
	}
	tagservice := tagservice.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	total, err := tagservice.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_TAG_FAIL, data)
		return
	}
	list, err := tagservice.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, data)
		return
	}
	data["lists"] = list
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTags(c *gin.Context) {
	appG := app.Gin{c}
	data := make(map[string]interface{})
	name := c.Query("name")
	state, _ := strconv.Atoi(c.DefaultQuery("state", "0"))
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, data)
		return
	}

	tagservice := tagservice.Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	exist, err := tagservice.ExistByName()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, data)
		return
	}
	if exist {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, data)
		return
	}
	tagservice.Add()
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTags(c *gin.Context) {
	appG := app.Gin{c}
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, data)
		return
	}

	tagservice := tagservice.Tag{
		ID:         id,
		Name:       name,
		ModifiedBy: modifiedBy,
		State:      state,
	}

	exist, err := tagservice.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, data)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, data)
		return
	}

	err = tagservice.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_TAG_FAIL, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} string "ok"
// @Failure 500 {object} string "fail"
// @Router /api/v1/tags/{id} [delete]
func DeleteTags(c *gin.Context) {
	appG := app.Gin{c}
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(c.Param("id"))
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, data)
		return
	}

	tagservice := tagservice.Tag{ID: id}
	exist, err := tagservice.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, data)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, data)
		return
	}
	err = tagservice.Delete()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, data)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

func ExportTag(c *gin.Context) {
	appG := app.Gin{c}
	name := c.PostForm("name")
	// fmt.Println(name)
	var state = -1
	if arg := c.PostForm("state"); arg != "" {
		state, _ = strconv.Atoi(arg)
	}
	tagservice := tagservice.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	filename, err := tagservice.Export()
	// fmt.Println(filename)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelFullPath() + filename,
	})

}
