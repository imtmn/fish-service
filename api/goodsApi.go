package api

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

type (
	// 商品分组
	GoodsGroup struct {
		Id        string         `json:"id"`
		Name      string         `json:"name"`
		Order     float64        `bson:"order" json:"order"`
		GoodsList []entity.Goods `json:"goods"`
	}
)

func SaveGoods(c *gin.Context) {
	goods := entity.Goods{}
	err := c.ShouldBindJSON(&goods)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	if goods.Id == "" {
		goods.StoreId = storeId
		// 默认：上架状态
		goods.Status = "on"
		goods.CreateTime = time.Now()
	}
	if goods.Id != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: goods.Name},
				{Key: "remark", Value: goods.Remark},
				{Key: "price", Value: goods.Price},
				{Key: "goodsType", Value: goods.GoodsType},
				{Key: "specification", Value: goods.Specification},
				{Key: "imgs", Value: goods.Imgs},
				{Key: "mainImags", Value: goods.MainImags},
				{Key: "updateTime", Value: time.Now()}}},
		}
		result := common.DB().UpdateByID("goods", goods.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("goods", goods)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}

}

func GetGoods(c *gin.Context) {
	goods := entity.Goods{}
	err := common.DB().FindID("goods", c.Param("id")).Decode(&goods)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对应的商品"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": goods})
}

func CountGoods(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "无店铺信息", "data": 0})
	}
	count := common.DB().CountByFilter("goods", bson.M{"storeId": storeId})
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": count})
}

func FindGoodsGroupByType(c *gin.Context) {
	// 默认最多查1000条数据
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "1000"), 10, 64)
	status := c.Query("status")
	storeId := c.Query("storeId")
	if storeId == "" {
		storeId = c.MustGet("storeId").(string)
	}
	filter := bson.M{
		"storeId": storeId,
	}
	if status != "" {
		filter["status"] = status
	}
	cursor := common.DB().CollectionDocumentsFilter("goods", 0, size, filter, bson.D{primitive.E{Key: "_id", Value: -1}})
	total := common.DB().CollectionCount("goods")
	goodsList := []entity.Goods{}
	cursor.All(context.TODO(), &goodsList)

	goodsGroup := []GoodsGroup{}

	for _, goods := range goodsList {
		goodsTypeList := goods.GoodsType
		if len(goodsTypeList) > 0 {
			// 有商品类型数据根据商品类型进行分组
			for _, goodsType := range goodsTypeList {
				addGoodsToGroupByType(&goodsGroup, goodsType, goods)
			}
		} else {
			// 没有分组的数据 放到其他类型中
			goodsType := entity.GoodsType{Name: "其他", Id: "-1"}
			addGoodsToGroupByType(&goodsGroup, goodsType, goods)
		}
	}

	typeCursor := common.DB().CollectionDocuments("goodsType", 0, size, bson.D{primitive.E{Key: "order", Value: 1}})
	goodstypeArr := []entity.GoodsType{}
	typeCursor.All(context.TODO(), &goodstypeArr)

	for i, _ := range goodsGroup {
		for _, gt := range goodstypeArr {
			if gt.Id == (&goodsGroup[i]).Id {
				(&goodsGroup[i]).Order = gt.Order
				break
			}
		}
		if (&goodsGroup[i]).Order == 0 {
			(&goodsGroup[i]).Order = 999999999999999
		}
	}

	sort.SliceStable(goodsGroup, func(i, j int) bool {
		return goodsGroup[i].Order < goodsGroup[j].Order
	})

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": goodsGroup, "total": total})
}

// 根据商品类型进行分组
func addGoodsToGroupByType(goodsGroup *[]GoodsGroup, goodsType entity.GoodsType, goods entity.Goods) {
	isExist := false
	var group *GoodsGroup
	// 遍历分组
	for i := 0; i < len(*goodsGroup); i++ {
		if (*goodsGroup)[i].Name == goodsType.Name {
			(*goodsGroup)[i].GoodsList = append((*goodsGroup)[i].GoodsList, goods)
			isExist = true
			break
		}
	}
	if !isExist {
		// 如果分组不存在 则新增一个分组
		group = &GoodsGroup{}
		group.Id = goodsType.Id
		group.Name = goodsType.Name
		//为分组添加元素
		group.GoodsList = append(group.GoodsList, goods)
		*goodsGroup = append(*goodsGroup, *group)
	}
}

func FindGoodsPage(c *gin.Context) {
	log.Println("page=", c.Query("page"), ",size=", c.Query("size"))
	page, _ := strconv.ParseInt(c.Param("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Param("size"), 10, 64)
	skip := (page - 1) * size
	cursor := common.DB().CollectionDocuments("goods", skip, size, bson.D{primitive.E{Key: "_id", Value: -1}})
	total := common.DB().CollectionCount("goods")
	goods := []entity.Goods{}
	cursor.All(context.TODO(), &goods)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取数据成功", "data": goods, "total": total})
}

func DeleteGoods(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	goods := entity.Goods{}
	err := common.DB().FindID("goods", c.Param("id")).Decode(&goods)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未找到对应的商品"})
		return
	}
	if goods.StoreId != storeId {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无数据删除权限"})
		return
	}
	id := c.Param("id")
	count := common.DB().DeleteById("goods", id)
	if count <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "删除数据失败", "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除数据成功", "count": count})
}

func ChangeStatusOn(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	id := c.Param("id")
	result := changeGoodsStatus(id, "on", storeId)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上架成功", "id": result.UpsertedID})
}

func ChangeStatusOff(c *gin.Context) {
	storeId := c.MustGet("storeId").(string)
	if storeId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "暂无店铺权限"})
		return
	}
	id := c.Param("id")
	result := changeGoodsStatus(id, "off", storeId)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "下架成功", "id": result.UpsertedID})
}

func changeGoodsStatus(id string, status string, storeId string) *mongo.UpdateResult {
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: status}, {Key: "updateTime", Value: time.Now()}}},
	}
	oid, _ := primitive.ObjectIDFromHex(id)
	return common.DB().UpdateOne("goods", bson.M{
		"_id":     oid,
		"storeId": storeId,
	}, update)
}
