package entity

import "time"

type (

	// 商品信息
	Goods struct {
		Id            string          `bson:"_id,omitempty" json:"id"`
		Name          string          `bson:"name" binding:"required" json:"name"`
		Remark        string          `bson:"remark" json:"remark"`
		Price         float64         `bson:"price" binding:"gte=1,lte=9999" json:"price,string"`
		Specification []Specification `bson:"specification" json:"specification"`
		GoodsType     []GoodsType     `bson:"goodsType" json:"goodsType"`
		Imgs          []string        `bson:"imgs" json:"imgs"`
		MainImags     string          `bson:"mainImags" json:"mainImags"`
		// on：已上架 off：已下架
		Status string `bson:"status" json:"status"`
		// 店铺Id
		StoreId    string    `bson:"storeId" json:"storeId"`
		Creator    string    `bson:"creator" json:"creator"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}

	// 商品类型
	GoodsType struct {
		Id    string  `bson:"_id,omitempty" json:"id" `
		Name  string  `bson:"name" binding:"required" json:"name"`
		Order float64 `bson:"order" json:"order"`
		// 店铺Id
		StoreId    string    `bson:"storeId" json:"storeId"`
		Creator    string    `bson:"creator" json:"creator"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
	}

	// 规格
	Specification struct {
		Id        string `bson:"_id,omitempty" json:"id" `
		Name      string `bson:"name" binding:"required" json:"name"`
		Price     string `bson:"price"  json:"price"`
		GroupName string `bson:"groupName"  json:"groupName"`
		// 店铺Id
		StoreId    string    `bson:"storeId" json:"storeId"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
	}
)
