package middleware

import (
	"library-api-user/internal/commons/response"
	"library-api-user/pkg/token"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedErrorWithAdditionalInfo(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}

func CheckAuthIsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedErrorWithAdditionalInfo(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		if payload.Role != "admin" {
			resp := response.UnauthorizedErrorWithAdditionalInfo("user doesn't have permission to access")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}

func CheckAuthIsAdminOrAuthor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.UnauthorizedErrorWithAdditionalInfo("len token must be 2")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := token.ValidateToken(bearerToken[1])
		if err != nil {
			resp := response.UnauthorizedErrorWithAdditionalInfo(err.Error())
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		if payload.Role == "user" {
			resp := response.UnauthorizedErrorWithAdditionalInfo("user doesn't have permission to access")
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}
		ctx.Set("authId", payload.AuthId)
		ctx.Set("role", payload.Role)
		ctx.Next()
	}
}
