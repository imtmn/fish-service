package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/medivhzhan/weapp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
)

//后台事件处理方法
func GetTokenByWxCode(code string) (string, error) {
	log.Println("code = ", code)
	loginResponse, err := OpenIDResolver(code)
	if err != nil {
		fmt.Println("获取不到openID,code=" + code)
		return "", fmt.Errorf("获取不到openID")
	}
	// //将openID存入数据库，返回对应_id
	user := &entity.User{}
	common.DB().FindOne("user", "wxOpenid", loginResponse.OpenID).Decode(user)
	if user.Id == "" {
		//数据库没有找到对应数据，应该将用户信息存储到数据库
		user.WxOpenid = loginResponse.OpenID
		user.Role = "user" // 默认为普通用户
		user.CreateTime = time.Now()
		user.UpdateTime = time.Now()
		result := common.DB().InsertOne("user", user)
		if result.InsertedID != nil {
			user.Id = result.InsertedID.(primitive.ObjectID).Hex()
		}
	}

	//使用accountID生成token
	userClaims := &UserClaims{
		UserId:   user.Id,
		OpenID:   user.WxOpenid,
		NickName: user.NickName,
		Role:     user.Role,
		StoreId:  user.StoreId,
		// 微信返回的sessionID
		SessionKey: loginResponse.SessionKey,
	}
	token, err := Sign(userClaims, 3600)
	if err != nil {
		return "", fmt.Errorf("不能生成token")
	}
	return token, nil
}

//将客户端上传的code，和小程序ID和秘钥上传至微信api换取openID
func OpenIDResolver(code string) (*weapp.LoginResponse, error) {
	appId := os.Getenv("AppId")
	appsecret := os.Getenv("Appsecret")
	resp, err := weapp.Login(appId, appsecret, code)
	if err != nil {
		return &weapp.LoginResponse{}, fmt.Errorf("weapp login: %v", err)
	}
	return &resp, nil
}
