package routers

import (
	"github.com/gin-gonic/gin"
)

// 注册路由
func RegisterRouters(r *gin.Engine) {

	apiGroup := r.Group("/api")
	loadV1(apiGroup)
}
