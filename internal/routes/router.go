package routes

import (
	"fmt"
	"library-api-user/internal/factory"
	"library-api-user/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(provider *factory.Provider) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger(), CORS())

	router.GET("/", func(ctx *gin.Context) {
		currentYear := time.Now().Year()
		message := fmt.Sprintf("Library API User %d", currentYear)

		ctx.JSON(http.StatusOK, message)
	})

	api := router.Group("/api")
	{
		v1 := api.Group("v1")
		{

			v1.POST("/register", provider.AuthProvider.Register)
			v1.POST("/login", provider.AuthProvider.Login)

			admin := v1.Use(middleware.CheckAuthIsAdmin())
			admin.GET("/users", provider.UserProvider.GetAll)
			admin.GET("/users/:id", provider.UserProvider.Detail)
			admin.PUT("/users/manage", provider.UserProvider.Update)

			auth := v1.Use(middleware.CheckAuth())
			auth.POST("/users/:id/borrow", provider.UserProvider.BorrowBook)
			auth.POST("/users/:id/return", provider.UserProvider.ReturnBook)
		}
	}

	return router
}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, accept, access-control-allow-origin, access-control-allow-headers")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}
