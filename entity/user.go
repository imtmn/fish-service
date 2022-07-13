package entity

import "time"

// 用户信息
type User struct {
	Id string `bson:"_id,omitempty" json:"id"`
	// 微信平台 openId
	WxOpenid string `bson:"wxOpenid" json:"wxOpenid"`
	// 角色：admin 超级管理员 boss 店长 manager 管理员 staff 店员 user 普通用户
	Role string `bson:"role" json:"role"`
	// 店铺Id
	StoreId string `bson:"storeId" json:"storeId"`
	// 网名
	NickName string `bson:"nickName" json:"nickName"`
	// 头像地址
	AvatarUrl string `bson:"avatarUrl" json:"avatarUrl"`
	// 手机信息
	Phone string `bson:"phone" json:"phone"`
	// 性别 男 女
	Gender string `bson:"gender" json:"gender"`
	// 地址信息
	Address    string    `bson:"address" json:"address"`
	CreateTime time.Time `bson:"createTime" json:"createTime"`
	UpdateTime time.Time `bson:"updateTime" json:"updateTime"`
}
