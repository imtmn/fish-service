package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/entity"
	"mtmn.top/fish-service/service"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("请求链接", c.Request.URL)
		jwtToken, error := service.VerifyTokenHeader(c)
		if error != nil {
			c.Abort()
			log.Println("未授权访问连接", c.Request.URL)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
			return
		}
		// log.Println("accessToken.Claims", accessToken.Claims.(jwt.MapClaims))
		claims := jwtToken.Claims.(*service.UserClaims)
		// 验证通过，会继续访问下一个中间件
		c.Set("claims", claims)
		c.Set("userId", claims.UserId)
		c.Set("storeId", claims.StoreId)
		c.Set("sessionKey", claims.SessionKey)
		c.Next()
	}
}

func DevAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		devUserId := os.Getenv("DEV_USER_ID")
		userInfo := entity.User{}
		log.Println(devUserId)
		common.DB().FindID("user", devUserId).Decode(&userInfo)
		log.Println(userInfo)
		if userInfo.Id == "" {
			ctx.Abort()
			log.Println("未授权访问连接", ctx.Request.URL)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
			return
		}

		userClaims := &service.UserClaims{
			UserId:   userInfo.Id,
			OpenID:   userInfo.WxOpenid,
			NickName: userInfo.NickName,
			Role:     userInfo.Role,
			StoreId:  userInfo.StoreId,
		}
		ctx.Set("claims", userClaims)
		ctx.Set("userId", userInfo.Id)
		ctx.Set("storeId", userInfo.StoreId)
		ctx.Next()

	}
}
