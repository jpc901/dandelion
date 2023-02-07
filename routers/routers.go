package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)

	// 登录业务路由
	r.POST("/login", controller.LoginHandler)

	r.GET("/", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是登录的用户, 判断请求头中是否有有效的JWT
		c.Request.Header.Get("Authorization")
		c.String(http.StatusOK, "登录成功")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
