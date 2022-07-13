package entity

import "time"

type (
	Feedback struct {
		Id      string   `bson:"_id,omitempty" json:"id"`
		Content string   `bson:"content" binding:"required" json:"content"`
		Imgs    []string `bson:"imgs" json:"imgs"`
		// 姓名
		Name string `bson:"name" json:"name"`
		// 电话
		Phone      string    `bson:"phone" json:"phone"`
		CreatorId  string    `bson:"creatorId" json:"creatorId"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}
)
