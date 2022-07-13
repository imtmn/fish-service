package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
	"go.mongodb.org/mongo-driver/bson"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/dto"
	"mtmn.top/fish-service/service"
)

func LoginByWx(c *gin.Context) {
	code := c.Query("code")
	token, err := service.GetTokenByWxCode(code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "用户登录失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "token": token})
}

// 解密微信用户基本信息
func DecodeWxUserInfo(c *gin.Context) {
	wxData := dto.WxData{}
	err := c.ShouldBindJSON(&wxData)
	ssk := c.MustGet("sessionKey").(string)
	userId := c.MustGet("userId").(string)
	log.Println("sessionKey = ", ssk)
	log.Println("userId = ", userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "获取参数失败"})
		return
	}
	wxUserInfo, err := weapp.DecryptUserInfo(wxData.RawData, wxData.EncryptedData, wxData.Signature, wxData.Iv, ssk)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "解析信息失败"})
		return
	}

	if wxUserInfo.OpenID != "" {
		gender := "未知"
		//0：未知、1：男、2：女
		if wxUserInfo.Gender == 1 {
			gender = "男"
		} else if wxUserInfo.Gender == 2 {
			gender = "女"
		}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "nickName", Value: wxUserInfo.Nickname},
				{Key: "address", Value: wxUserInfo.Province + wxUserInfo.City + wxUserInfo.Country},
				{Key: "avatarUrl", Value: wxUserInfo.Avatar},
				{Key: "gender", Value: gender},
				{Key: "updateTime", Value: time.Now()},
			}},
		}
		common.DB().UpdateByID("user", userId, update)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "token": "解析成功"})
}

// 解密微信用户手机号
func DecodeWxPhone(c *gin.Context) {
	wxData := dto.WxData{}
	err := c.ShouldBindJSON(&wxData)
	ssk := c.MustGet("sessionKey").(string)
	userId := c.MustGet("userId").(string)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "获取参数失败"})
		return
	}
	phone, err := weapp.DecryptPhoneNumber(ssk, wxData.EncryptedData, wxData.Iv)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "解析信息失败"})
		return
	}

	if phone.PhoneNumber != "" {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "phone", Value: phone.PhoneNumber},
				{Key: "updateTime", Value: time.Now()},
			}},
		}
		common.DB().UpdateByID("user", userId, update)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "token": "解析成功"})
}
