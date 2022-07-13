package entity

import "time"

type (

	// 订单
	Order struct {
		Id string `bson:"_id,omitempty" json:"id"`
		// 备注
		Remark string `bson:"remark" json:"remark"`
		// 订单金额
		Price float64 `bson:"price" json:"price,string"  binding:"gte=1,lte=9999"`
		// 支付金额
		PayPrice float64 `bson:"payPrice" json:"payPrice,string"`
		// normal：待支付 pay：已支付 cancel:取消 complete:完成
		Status string `bson:"status" json:"status"`
		// 店铺Id
		StoreId string `bson:"storeId" json:"storeId"`
		// 店铺名称
		StoreName string `bson:"storeName" json:"storeName"`
		// 桌号iD
		DeskId string `bson:"deskId" json:"deskId"`
		// 桌号名称
		DeskName string `bson:"deskName" json:"deskName"`
		// 人数
		PersonNum int16 `bson:"personNum" json:"personNum" binding:"gte=1,lte=99"`
		// 创建人ID
		CreatorId string `bson:"creatorId" json:"creatorId"`
		// 订单商品信息
		OrderGoods []OrderGoods `bson:"orderGoods" json:"orderGoods"`
		CreateTime time.Time    `bson:"createTime" json:"createTime"`
		UpdateTime time.Time    `bson:"updateTime" json:"updateTime"`
	}

	// 订单商品
	OrderGoods struct {
		Id     string  `bson:"_id,omitempty" json:"id"`
		Name   string  `bson:"name" json:"name"`
		Price  float64 `bson:"price" binding:"gte=1,lte=9999" json:"price,string"`
		Number int8    `bson:"number" binding:"gte=1,lte=999" json:"number"`
		// 店铺Id
		StoreId    string    `bson:"storeId" json:"storeId"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}
)
