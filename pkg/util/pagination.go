package util

import (
	"gin/pkg/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(ctx *gin.Context) int {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page <= 0 {
		return 0
	}
	return (page - 1) * setting.PageSize
}
