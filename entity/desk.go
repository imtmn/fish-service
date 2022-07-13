package entity

import "time"

type (
	Desk struct {
		Id   string `bson:"_id,omitempty" json:"id"`
		Name string `bson:"name" binding:"required" json:"name"`
		// 店铺Id
		StoreId    string    `bson:"storeId" json:"storeId"`
		CreatorId  string    `bson:"creatorId" json:"creatorId"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}
)
