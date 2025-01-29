package controllers

import (
	"library-api-user/internal/commons/response"
	"library-api-user/internal/models"
	"library-api-user/internal/params"
	"library-api-user/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Detail(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	BorrowBook(ctx *gin.Context)
	ReturnBook(ctx *gin.Context)
}

type UserControllerImpl struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Detail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	result, custErr := controller.UserService.Detail(ctx, uint64(id))
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success retrieve data detail user", result)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *UserControllerImpl) Update(ctx *gin.Context) {
	var req = new(params.UserRequest)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	role := ctx.GetString("role")
	if role != "admin" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "user doesn't have permission to manage",
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}
	req.ID = uint64(id)

	authId := ctx.GetInt("authId")

	custErr := controller.UserService.Update(ctx, req, uint64(authId))

	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccess()
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *UserControllerImpl) GetAll(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	pageNum := 1
	limitSize := 5

	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			pageNum = parsedPage
		}
	}

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 {
			limitSize = parsedLimit
		}
	}

	pagination := models.Pagination{
		Page:     pageNum,
		Offset:   (pageNum - 1) * limitSize,
		PageSize: limitSize,
	}

	result, custErr := controller.UserService.GetAll(ctx, &pagination)

	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	type Response struct {
		Users      interface{} `json:"users"`
		Pagination interface{} `json:"pagination"`
	}

	var responses Response
	responses.Users = result
	responses.Pagination = pagination

	resp := response.GeneralSuccessCustomMessageAndPayload("Success get data users", responses)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *UserControllerImpl) BorrowBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	authID := ctx.GetInt("authId")

	custErr := controller.UserService.BorrowBook(ctx, uint64(authID), uint64(id))
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success users borrow book", nil)
	ctx.JSON(resp.StatusCode, resp)
}

func (controller *UserControllerImpl) ReturnBook(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err,
		})
		return
	}

	authId := ctx.GetInt("authId")

	custErr := controller.UserService.ReturnBook(ctx, uint64(authId), uint64(id))
	if custErr != nil {
		ctx.AbortWithStatusJSON(custErr.StatusCode, custErr)
		return
	}

	resp := response.GeneralSuccessCustomMessageAndPayload("Success users return book", nil)
	ctx.JSON(resp.StatusCode, resp)
}
