package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

func SaveGoodsType(c *gin.Context) {
	goodsType := entity.GoodsType{}
	err := c.ShouldBindJSON(&goodsType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	storeId := c.MustGet("storeId").(string)

	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}

	if goodsType.Id == "" {
		goodsType.StoreId = storeId
		goodsType.CreateTime = time.Now()
		goodsType.Order = float64(time.Now().Unix())
	}
	if goodsType.Id != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: goodsType.Name},
				{Key: "order", Value: goodsType.Order},
			}},
		}
		result := common.DB().UpdateByID("goodsType", goodsType.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("goodsType", goodsType)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}
}

func GetGoodsType(c *gin.Context) {
	goodsType := entity.GoodsType{}
	err := common.DB().FindID("goodsType", c.Param("id")).Decode(&goodsType)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对应商品类型"})
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": goodsType})
}

func FindGoodsTypePage(c *gin.Context) {
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "100"), 10, 64)
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	filter := bson.M{
		"storeId": storeId,
	}
	cursor := common.DB().CollectionDocumentsFilter("goodsType", 0, size, filter, bson.D{primitive.E{Key: "order", Value: 1}})
	total := common.DB().CollectionCount("goodsType")
	goodstype := []entity.GoodsType{}
	cursor.All(context.TODO(), &goodstype)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": goodstype, "total": total})
}

func DeleteGoodsType(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	id := c.Param("id")
	goodsType := entity.GoodsType{}
	err := common.DB().FindID("goodsType", id).Decode(&goodsType)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对应商品类型"})
	}
	if goodsType.StoreId != storeId {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无数据删除权限"})
		return
	}
	count := common.DB().DeleteById("goodsType", id)
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "删除数据失败", "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除数据成功", "count": count})
}
