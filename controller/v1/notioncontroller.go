package v1

import (
	"github.com/gin-gonic/gin"
	"notion-native/model/request"
	"notion-native/service"
)

func CalculateHoursPost(r *gin.Context) {
	caclulateHours := &request.CalculateHours{}
	if err := r.ShouldBindJSON(caclulateHours); err != nil {
		r.JSON(200, gin.H{
			"message": err.Error(),
			"status":  400,
		})
		return
	}
	total, err := service.CalculateHours(caclulateHours.PageId, caclulateHours.OrderId, caclulateHours.NotionVersion)
	data := make(map[string]interface{})
	data["total"] = total
	if err != nil {
		r.JSON(200, gin.H{
			"status":  400,
			"message": "系统错误",
		})
	} else {
		r.JSON(200, gin.H{
			"status":  200,
			"message": "",
			"data":    data,
		})
	}
}

func CalculateHoursGet(r *gin.Context) {
	pageId := r.Params.ByName("pageId")
	orderId := r.Params.ByName("orderId")
	notionVersion := r.Params.ByName("notionVersion")
	if pageId == "" || orderId == "" || notionVersion == "" {
		r.JSON(200, gin.H{
			"status":  400,
			"message": "参数错误",
		})
		return
	}
	total, err := service.CalculateHours(pageId, orderId, notionVersion)
	data := make(map[string]interface{})
	data["total"] = total
	if err != nil {
		r.JSON(200, gin.H{
			"status":  400,
			"message": "系统错误",
		})
	} else {
		r.JSON(200, gin.H{
			"status":  200,
			"message": "",
			"data":    data,
		})
	}
}
