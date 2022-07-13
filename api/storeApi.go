package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
	"mtmn.top/fish-service/service"
)

// 保存店铺申请
func SaveApplicationForm(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	form := &entity.ApplicationForm{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if form.Id == "" {
		// 基础信息
		form.Status = "apply" //申请中
		form.CreatorId = userId
		form.CreateTime = time.Now()
		form.UpdateTime = time.Now()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "您已经提交申请请求，等待审核，请勿重复提交"})
		return
	}
	result := common.DB().InsertOne("applicationForm", form)
	if result.InsertedID != nil {
		applyId := result.InsertedID.(primitive.ObjectID).Hex()
		form.Id = applyId
		go common.SendMailApply(form)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
}

// 获取申请
func GetApplicationForm(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	if userId != "" {
		form := &entity.ApplicationForm{}
		common.DB().FindOne("applicationForm", "creatorId", userId).Decode(&form)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": form})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "无申请数据"})
}

func SaveStore(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	store := entity.Store{}
	err := c.ShouldBindJSON(&store)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if store.Id == "" {
		// 店铺基础信息
		store.Creator = userId
		store.CreateTime = time.Now()
		store.UpdateTime = time.Now()
	}
	log.Println(store)
	if store.Id != "" {
		//FIXME 店铺管理员才有权限修改店铺信息
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: store.Name},
				{Key: "remark", Value: store.Remark},
				{Key: "avatarPath", Value: store.AvatarPath},
				{Key: "hidePrice", Value: store.HidePrice},
				{Key: "updateTime", Value: store.UpdateTime},
			}},
		}
		result := common.DB().UpdateByID("store", store.Id, update)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.UpsertedID})
	} else {
		result := common.DB().InsertOne("store", store)
		if result.InsertedID != nil {
			storeId := result.InsertedID.(primitive.ObjectID).Hex()
			log.Println("创建店铺成功，店铺id:", storeId, "用户id:"+userId)
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "role", Value: "boss"},
					{Key: "storeId", Value: storeId},
				}},
			}
			common.DB().UpdateByID("user", userId, update)
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "保存成功", "id": result.InsertedID})
	}
}

// 如果有ID 则根据店铺ID 获取店铺信息 如果参数没有ID 是否有创建店铺
func GetStore(c *gin.Context) {
	storeId := c.Param("id")
	if storeId == "" {
		storeId = c.MustGet("storeId").(string)
	}
	if storeId != "" {
		store := &entity.Store{}
		common.DB().FindID("store", storeId).Decode(store)
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": store})
	}
}

// 生成邀请码
func GenerateInviteCode(c *gin.Context) {
	//  只能生成一个邀请连接 -> redis key：店铺主键 value:生成的唯一编码，过期时间设置为10分钟
	user := c.MustGet("claims").(*service.UserClaims)
	storeId := user.StoreId
	if storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "获取店铺信息失败，无法生成邀请链接，请确认您是否已经创建店铺成功"})
		return
	}
	uuid := uuid.New().String()
	err := common.RedisClient().Set(context.Background(), storeId, uuid, time.Minute*10).Err()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "生成邀请码失败"})
		return
	}
	log.Println("生成店铺邀请码，", "店铺ID为：", storeId, "邀请编码为：", uuid)
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取成功", "data": uuid})
}

// 根据邀请码加入店铺
func JoinStore(c *gin.Context) {
	code := c.Query("code")
	storeId := c.Query("storeId")
	userId := c.MustGet("userId").(string)
	if code == "" || storeId == "" {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "加入店铺失败，参数有误"})
	}
	log.Println("通过邀请码加入店铺：", "店铺ID为：", storeId, "邀请编码为：", code)
	stringCmd := common.RedisClient().Get(context.Background(), storeId)
	if stringCmd.Err() != nil {
		log.Println(stringCmd.Err())
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "加入店铺失败，邀请码获取失败，请尝试重新邀请"})
		return
	}
	// 根据店铺ID redis 验证编码是否存在，如果编码不存在，提醒前端连接已经过期，需要进行重新邀请
	if stringCmd.Val() != code {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "加入店铺失败，邀请码已过期"})
		return
	}
	user := entity.User{}
	err := common.DB().FindID("user", userId).Decode(&user)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "获取用户信息失败"})
		return
	}
	// 如果编码存在，校验该用户是否已经加入其他店铺，如果已经加入，前端提醒不允许再加入
	if user.StoreId != "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "您已经加入店铺，请退出当前店铺再接受邀请"})
		return
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "storeId", Value: storeId},
			{Key: "role", Value: "staff"},
			{Key: "updateTime", Value: time.Now()},
		}},
	}
	result := common.DB().UpdateByID("user", userId, update)

	if result.UpsertedCount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "个人信息异常，加入店铺失败"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"code": 200, "msg": "加入店铺成功"})
}

// 店铺申请审批通过
func ApplyApss(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 200, "msg": "id为空"})
	}
	applyForm := &entity.ApplicationForm{}
	common.DB().FindID("applicationForm", id).Decode(applyForm)
	if applyForm.Status != "apply" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "为找到对应申请信息"})
		return
	}
	store := &entity.Store{}
	store.Name = applyForm.Name
	store.Remark = applyForm.Remark
	store.Creator = applyForm.CreatorId
	store.AvatarPath = "/upload/202207/1657237593-07FzCle0wdp70ac5e049b118708bd43f057e0567eb2a.jpg"
	store.CreateTime = time.Now()
	store.UpdateTime = time.Now()
	insertOneResult := common.DB().InsertOne("store", store)
	common.DB().UpdateByID("user", applyForm.CreatorId, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "storeId", Value: insertOneResult.InsertedID.(primitive.ObjectID).Hex()},
			{Key: "updateTime", Value: time.Now()},
		}},
	})
	common.DB().UpdateByID("applicationForm", id, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: "pass"},
			{Key: "updateTime", Value: time.Now()},
		}},
	})
	go InitStoreData(applyForm.CreatorId, insertOneResult.InsertedID.(primitive.ObjectID).Hex())
}

// 初始化店铺数据
func InitStoreData(userId, storeId string) {
	// 初始化商品类型
	goodsType := entity.GoodsType{}
	goodsType.Name = "示例类型"
	goodsType.CreateTime = time.Now()
	goodsType.StoreId = storeId
	goodsType.Creator = userId
	typeResult := common.DB().InsertOne("goodsType", goodsType)
	goodsType.Id = typeResult.InsertedID.(primitive.ObjectID).Hex()
	// 初始化商品
	var typeArr = []entity.GoodsType{goodsType}
	goods := entity.Goods{}
	goods.Name = "示例菜品"
	goods.Remark = "美味佳肴，进店必点"
	goods.Status = "on"
	goods.MainImags = "/upload/202207/1657237191-bBEpch5HkzN78ccb0f9c179320f8c4b5a8f6547d765c.jpg"
	goods.CreateTime = time.Now()
	goods.StoreId = storeId
	goods.Creator = userId
	goods.GoodsType = typeArr
	common.DB().InsertOne("goods", goods)
	// 初始化桌号
	desk := entity.Desk{}
	desk.Name = "示例桌号"
	desk.CreateTime = time.Now()
	desk.CreatorId = userId
	desk.StoreId = storeId
	common.DB().InsertOne("desk", desk)
}

// 店铺申请 审批不通过
func ApplyReject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 200, "msg": "id为空"})
	}
	applyForm := &entity.ApplicationForm{}
	common.DB().FindID("applicationForm", id).Decode(applyForm)
	if applyForm.Status != "apply" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "为找到对应申请信息"})
		return
	}
	common.DB().UpdateByID("applicationForm", id, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: "reject"},
			{Key: "updateTime", Value: time.Now()},
		}},
	})
}
