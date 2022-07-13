package entity

import "time"

type (
	// 店铺信息
	Store struct {
		Id         string `bson:"_id,omitempty" json:"id"`
		Name       string `bson:"name" binding:"required" json:"name"`
		Remark     string `bson:"remark" json:"remark"`
		AvatarPath string `bson:"avatarPath" json:"avatarPath"`
		// 如果为true,界面上显示价格
		HidePrice  bool      `bson:"hidePrice" json:"hidePrice"`
		Creator    string    `bson:"creator" json:"creator"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}

	// 申请表
	ApplicationForm struct {
		Id string `bson:"_id,omitempty" json:"id"`
		// 店铺名称
		Name string `bson:"name" binding:"required" json:"name"`
		// 店铺简介
		Remark string `bson:"remark"  json:"remark"`
		// 附件列表
		Imgs []string `bson:"imgs" json:"imgs"`
		// 姓名
		PersonName string `bson:"personName" json:"personName"`
		// 电话
		Phone string `bson:"phone" json:"phone"`
		// 补充说明
		Info string `bson:"info"  json:"info"`
		// 审批状态 ： apply：申请中 pass：审批通过 reject：审批不通过
		Status string `bson:"status"  json:"status"`
		// 审批结果说明
		Result     string    `bson:"result"  json:"result"`
		CreatorId  string    `bson:"creatorId" json:"creatorId"`
		CreateTime time.Time `bson:"createTime" json:"createTime"`
		UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
	}
)
