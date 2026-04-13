package v1

import (
	"github.com/gin-gonic/gin"
	"notion-native/model/request"
	"notion-native/service"
)

func CalculateHours(r *gin.Context) {
	caclulateHours := &request.CalculateHours{}
	if err := r.ShouldBindJSON(caclulateHours); err != nil {
		r.JSON(200, gin.H{
			"message": err.Error(),
			"status":  401,
		})
		return
	}
	total, err := service.CalculateHours(caclulateHours.PageId, caclulateHours.OrderId, caclulateHours.NotionVersion)
	data := make(map[string]interface{})
	data["total"] = total
	if err != nil {
		r.JSON(200, gin.H{
			"status":  401,
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
