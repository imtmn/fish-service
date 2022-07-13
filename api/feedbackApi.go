package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

func SaveFeedBack(c *gin.Context) {
	feedback := entity.Feedback{}
	err := c.ShouldBindJSON(&feedback)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	userId := c.MustGet("userId").(string)
	if feedback.Id == "" {
		feedback.CreateTime = time.Now()
		feedback.UpdateTime = time.Now()
		feedback.CreatorId = userId
	}
	result := common.DB().InsertOne("feedback", feedback)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
}
