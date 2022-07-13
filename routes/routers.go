package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"mtmn.top/fish-service/api"
	"mtmn.top/fish-service/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	dev := os.Getenv("DEV_USER_ID")
	v1 := router.Group("/api/v1")
	{
		if dev == "" {
			v1.Use(middleware.JwtAuth())
		} else {
			v1.Use(middleware.DevAuth())
		}

		// 文件管理
		v1.POST("/upload/token", api.OSSUploadFile)

		// 菜品管理
		v1.POST("/goods/save", api.SaveGoods)
		v1.GET("/goods/:id", api.GetGoods)
		// v1.GET("/goods", api.FindGoodsPage)
		v1.GET("/goods/groups", api.FindGoodsGroupByType)
		v1.DELETE("/goods/:id", api.DeleteGoods)
		// 下架商品
		v1.GET("/goods/off/:id", api.ChangeStatusOff)
		// 上架商品
		v1.GET("/goods/on/:id", api.ChangeStatusOn)
		//菜品数
		v1.GET("/goods/count", api.CountGoods)

		// 品类管理
		v1.POST("/goodsType/save", api.SaveGoodsType)
		v1.GET("/goodsType/:id", api.GetGoodsType)
		v1.GET("/goodsType", api.FindGoodsTypePage)
		v1.DELETE("/goodsType/:id", api.DeleteGoodsType)

		// 规格管理
		v1.POST("/specification/save", api.SaveSpecification)
		v1.GET("/specification/:id", api.GetSpecification)
		v1.GET("/specification", api.FindSpecificationPage)
		v1.DELETE("/specification/:id", api.DeleteSpecification)

		// 订单管理
		v1.POST("/order/save", api.SaveOrder)
		v1.GET("/order/:id", api.GetOrder)
		v1.GET("/order", api.FindOrderPage)
		v1.GET("/order/cancel/:id", api.CancelOrder)
		v1.GET("/order/recover/:id", api.RecoverOrder)
		v1.GET("/order/pay", api.PayOrder)
		v1.GET("/order/count/today", api.CountTodayOrder)
		v1.GET("/order/income/today", api.IncomeOrder)

		// 店铺信息
		v1.POST("/store/save", api.SaveStore)
		v1.GET("/store/:id", api.GetStore)
		v1.GET("/store", api.GetStore)
		v1.GET("/store/code", api.GenerateInviteCode)
		v1.GET("/store/join", api.JoinStore)
		v1.POST("/store/apply", api.SaveApplicationForm)
		v1.GET("/store/apply", api.GetApplicationForm)

		// 用户信息
		v1.POST("/user/save", api.SaveUser)
		v1.GET("/user/store", api.FindUserByStore)
		v1.GET("/user", api.GetUser)
		v1.DELETE("/user/remove/:userId", api.RemoveUserByStore)
		v1.GET("/user/role", api.ChangeUserRole)

		// 桌号管理
		v1.POST("/desk/save", api.SaveDesk)
		v1.GET("/desk", api.FindDeskPage)
		v1.GET("/desk/qrcode/:id", api.GetQrcode)
		v1.GET("/desk/:id", api.GetDesk)
		v1.DELETE("/desk/:id", api.DeleteDesk)

		// 问题反馈
		v1.POST("/feedback/save", api.SaveFeedBack)

	}

	common := router.Group("/api/v1")
	{
		//公共接口
		common.GET("/login", api.LoginByWx)
		// 审批接口 权限后面逐步完善
		common.GET("/apply/pass/:id", api.ApplyApss)
		common.GET("/apply/reject/:id", api.ApplyReject)
	}

	return router
}
