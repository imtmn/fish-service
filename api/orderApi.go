package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

// 保存订单
func SaveOrder(c *gin.Context) {
	order := entity.Order{}
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	goodsList := order.OrderGoods
	if len(goodsList) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "订单不能为空"})
		return
	}

	// 核对订单商品价格
	var totalPrice float64
	for _, orderGoods := range goodsList {
		goods := entity.Goods{}
		log.Println("开始获取商品")
		err := common.DB().FindID("goods", orderGoods.Id).Decode(&goods)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "获取商品信息失败，请重试"})
			return
		}
		if orderGoods.Price != goods.Price {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "订单商品价格有误"})
			return
		}
		totalPrice += orderGoods.Price * float64(orderGoods.Number)
	}

	// 确认商品订单总价
	if order.Price != totalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "订单商品价格有误"})
		return
	}
	if order.Id == "" {
		order.Status = "normal"
		order.CreatorId = c.MustGet("userId").(string)
		order.CreateTime = time.Now()
		order.UpdateTime = time.Now()
	}
	if order.Id != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "remark", Value: order.Remark}}},
		}
		result := common.DB().UpdateByID("order", order.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("order", order)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}
}

func GetOrder(c *gin.Context) {
	order := entity.Order{}
	log.Println("开始获取商品")
	storeId := c.MustGet("storeId").(string)
	userId := c.MustGet("userId").(string)
	err := common.DB().FindID("order", c.Param("id")).Decode(&order)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对于商品类型"})
	}
	if order.StoreId != storeId && order.CreatorId != userId {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "没有权限查看该订单"})
		return
	}
	log.Println(order.Id)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": order})
}

func FindOrderPage(c *gin.Context) {
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "100"), 10, 64)
	storeId := c.MustGet("storeId").(string)
	userId := c.MustGet("userId").(string)
	filter := bson.M{"$or": []bson.M{{"storeId": storeId}, {"creatorId": userId}}}
	if c.Query("status") != "" {
		filter["status"] = c.Query("status")
	}
	sort := bson.D{primitive.E{Key: "createTime", Value: -1}}
	cursor := common.DB().CollectionDocumentsFilter("order", 0, size, filter, sort)
	total := common.DB().CollectionCount("order")
	order := []entity.Order{}
	cursor.All(context.TODO(), &order)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": order, "total": total})
}

func CancelOrder(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	coll := common.DB().Coll("order")
	if oid, err := primitive.ObjectIDFromHex(c.Param("id")); err == nil {
		result, err := coll.UpdateOne(
			context.TODO(),
			bson.M{"_id": oid, "status": "normal", "storeId": storeId},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "status", Value: "cancel"}}},
			},
		)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "取消订单成功", "data": result.UpsertedID})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "取消订单失败，请刷新后重试"})
}

func RecoverOrder(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	coll := common.DB().Coll("order")
	if oid, err := primitive.ObjectIDFromHex(c.Param("id")); err == nil {
		result, err := coll.UpdateOne(
			context.TODO(),
			bson.M{"_id": oid, "status": "cancel", "storeId": storeId},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "status", Value: "normal"}}},
			},
		)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "恢复订单成功", "data": result.UpsertedID})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "恢复订单失败，请刷新后重试"})
}

func PayOrder(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	coll := common.DB().Coll("order")
	payPrice, err := strconv.ParseFloat(c.Query("payPrice"), 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "结算金额有误，请确认"})
		return
	}
	log.Println("支付金额", payPrice)
	log.Println("订单号", c.Query("id"))
	if oid, err := primitive.ObjectIDFromHex(c.Query("id")); err == nil {
		result, err := coll.UpdateOne(
			context.TODO(),
			bson.M{"_id": oid, "status": "normal", "storeId": storeId},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "status", Value: "pay"},
					{Key: "payPrice", Value: payPrice}}},
			},
		)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "结算成功", "data": result.UpsertedID})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "结算失败，请刷新后重试"})
}

func CountTodayOrder(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "无店铺信息", "data": 0})
	}
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	count := common.DB().CountByFilter("order", bson.M{
		"storeId":    storeId,
		"createTime": bson.M{"$gt": startTime},
	})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": count})
}

func IncomeOrder(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "无店铺信息", "data": 0})
	}
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "storeId",
				Value: bson.D{
					{Key: "$eq", Value: storeId},
				},
			},
			{Key: "createTime",
				Value: bson.D{
					{Key: "$gt", Value: startTime},
				},
			},
		}},
	}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "storeId", Value: "$storeId"},
			}},
			{Key: "sum", Value: bson.D{
				{Key: "$sum", Value: "$payPrice"},
			}},
		}},
	}
	cursor, err := common.DB().Coll("order").Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": 0})
	}
	for _, result := range results {
		fmt.Printf("Average price of %v \n", result["payPrice"])
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": results[0]["sum"]})
}
