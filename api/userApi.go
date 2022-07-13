package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
	"mtmn.top/fish-service/service"
)

// 保存用户信息
func SaveUser(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	user := entity.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "nickName", Value: user.NickName},
			{Key: "phone", Value: user.Phone},
			{Key: "gender", Value: user.Gender},
			{Key: "address", Value: user.Address},
			{Key: "avatarUrl", Value: user.AvatarUrl},
			{Key: "updateTime", Value: time.Now()},
		}},
	}
	result := common.DB().UpdateByID("user", userId, update)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
}

// 获取用户信息
func GetUser(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	user := entity.User{}
	err := common.DB().FindID("user", userId).Decode(&user)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取用户信息失败"})
	}
	log.Println(user.Id)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": user})
}

// 查询店铺的所有人员
func FindUserByStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	filter := bson.M{"storeId": storeId}
	sort := bson.D{primitive.E{Key: "createTime", Value: 1}}
	cursor := common.DB().CollectionDocumentsFilter("user", 0, 100, filter, sort)
	user := []entity.User{}
	cursor.All(context.TODO(), &user)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": user, "total": len(user)})
}

// 将人员从店铺中移除，店长才有权限
func RemoveUserByStore(c *gin.Context) {
	userClaims := c.MustGet("claims").(*service.UserClaims)
	if userClaims.Role != "boss" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "只有店长才有权限移除员工"})
	}
	storeId := c.MustGet("storeId").(string)
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "请选择需要移除的员工"})
	}
	oid, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"storeId": storeId, "_id": oid}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "storeId", Value: ""},
		}},
	}
	result := common.DB().UpdateOne("user", filter, update)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "人员移除成功", "id": result.UpsertedID})
}

// 修改用户权限
func ChangeUserRole(c *gin.Context) {
	userClaims := c.MustGet("claims").(*service.UserClaims)
	if userClaims.Role != "boss" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "只有店长才有权限移除员工"})
	}
	storeId := c.MustGet("storeId").(string)
	userId := c.Query("userId")
	role := c.Query("role")
	oid, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"storeId": storeId, "_id": oid}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "role", Value: role},
		}},
	}
	result := common.DB().UpdateOne("user", filter, update)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "修改成功", "id": result.UpsertedID})
}
