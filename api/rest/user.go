package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/dto"
	"github.com/m6yf/bcwork/utils"
)

// UserGetHandler Get users data.
// @Description Get users data.
// @Tags User
// @Param options body core.UserOptions true "Options"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.User
// @Security ApiKeyAuth
// @Router /user/get [post]
func (o *OMSNewPlatform) UserGetHandler(c *fiber.Ctx) error {
	data := &core.UserOptions{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for getting users data", err)
	}

	users, err := o.userService.GetUsers(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve users data", err)
	}

	return c.JSON(users)
}

// UserGetInfoHandler Get user info.
// @Description Get user info.
// @Tags User
// @Param id query string true "User ID"
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Security ApiKeyAuth
// @Router /user/info [get]
func (o *OMSNewPlatform) UserGetInfoHandler(c *fiber.Ctx) error {
	userID := c.Query("id")

	user, err := o.userService.GetUserInfo(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve user info", err)
	}

	return c.JSON(user)
}

// UserSetHandler Create new user.
// @Description Create new user.
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.User true "User"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /user/set [post]
func (o *OMSNewPlatform) UserSetHandler(c *fiber.Ctx) error {
	data := &dto.User{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for creating user", err)
	}

	err := o.userService.CreateUser(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create user", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user successfully created")
}

// UserUpdateHandler Update user.
// @Description Update user.
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.User true "User"
// @Success 200 {object} utils.BaseResponse
// @Security ApiKeyAuth
// @Router /user/update [post]
func (o *OMSNewPlatform) UserUpdateHandler(c *fiber.Ctx) error {
	data := &dto.User{}
	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "failed to parse request for updating user", err)
	}

	err := o.userService.UpdateUser(c.Context(), data)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update user", err)
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "user successfully updated")
}
