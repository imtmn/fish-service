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

func SaveSpecification(c *gin.Context) {
	specification := entity.Specification{}
	err := c.ShouldBindJSON(&specification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if specification.Id == "" {
		specification.CreateTime = time.Now()
	}
	log.Println(specification)
	if specification.Id != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: specification.Name}, {Key: "groupName", Value: specification.GroupName}}},
		}
		result := common.DB().UpdateByID("specification", specification.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("specification", specification)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}
}

func GetSpecification(c *gin.Context) {
	specification := entity.Specification{}
	log.Println("开始获取商品")
	err := common.DB().FindID("specification", c.Param("id")).Decode(&specification)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对应商品类型"})
	}
	log.Println(specification.Id)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": specification})
}

func FindSpecificationPage(c *gin.Context) {
	log.Println("page=", c.Query("page"), ",size=", c.Query("size"))
	page, _ := strconv.ParseInt(c.Param("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Param("size"), 10, 64)
	skip := (page - 1) * size
	cursor := common.DB().CollectionDocuments("specification", skip, size, bson.D{primitive.E{Key: "order", Value: 1}})
	total := common.DB().CollectionCount("specification")
	specification := []entity.Specification{}
	cursor.All(context.TODO(), &specification)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": specification, "total": total})
}

func DeleteSpecification(c *gin.Context) {
	id := c.Param("id")
	count := common.DB().DeleteById("specification", id)
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "删除数据失败", "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除数据成功", "count": count})
}
