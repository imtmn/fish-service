package api

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

func SaveDesk(c *gin.Context) {
	desk := entity.Desk{}
	err := c.ShouldBindJSON(&desk)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	storeId := c.MustGet("storeId").(string)
	userId := c.MustGet("userId").(string)
	if storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": "请先申请成为商家"})
		return
	}
	if desk.Id == "" {
		desk.CreateTime = time.Now()
		desk.UpdateTime = time.Now()
		desk.StoreId = storeId
		desk.CreatorId = userId
	}
	if desk.Id != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: desk.Name},
				{Key: "updateTime", Value: time.Now()},
			}},
		}
		result := common.DB().UpdateByID("desk", desk.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("desk", desk)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}
}

func GetDesk(c *gin.Context) {
	desk := entity.Desk{}
	err := common.DB().FindID("desk", c.Param("id")).Decode(&desk)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对桌号"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": desk})
}

func FindDeskPage(c *gin.Context) {
	log.Println("page=", c.Query("page"), ",size=", c.Query("size"))
	page, _ := strconv.ParseInt(c.Param("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Param("size"), 10, 64)
	skip := (page - 1) * size
	storeId := c.MustGet("storeId").(string)
	filter := bson.M{"storeId": storeId}
	cursor := common.DB().CollectionDocumentsFilter("desk", skip, size, filter, bson.D{primitive.E{Key: "createTime", Value: 1}})
	total := common.DB().CollectionCount("desk")
	desk := []entity.Desk{}
	cursor.All(context.TODO(), &desk)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": desk, "total": total})
}

func DeleteDesk(c *gin.Context) {
	id := c.Param("id")
	count := common.DB().DeleteById("desk", id)
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "删除数据失败", "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除数据成功", "count": count})
}

//获取二维码
func GetQrcode(c *gin.Context) {
	desk := entity.Desk{}
	err := common.DB().FindID("desk", c.Param("id")).Decode(&desk)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对桌号"})
		return
	}
	png, err := qrcode.Encode("https://mtmn.top/ucode/"+desk.Id, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "生成二维码失败"})
		return
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(png))
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "生成成功", "data": sEnc})

}
