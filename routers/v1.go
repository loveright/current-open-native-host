package routers

import (
	"github.com/gin-gonic/gin"
	v1 "notion-native/controller/v1"
)

func loadV1(g *gin.RouterGroup) {

	apiV1 := g.Group("/v1")

	notion := apiV1.Group("/notion")
	{
		notion.POST("/calculate/hours", v1.CalculateHoursPost)
		notion.GET("/calculate/hours/:pageId/:orderId/:notionVersion", v1.CalculateHoursGet)
	}
}
