package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/database"
	"github.com/muhafs/go-fiber-gorm/model/entity"
	"github.com/muhafs/go-fiber-gorm/model/request"
	"github.com/muhafs/go-fiber-gorm/utils"
	"gorm.io/gorm"
)

const (
	ErrorCredentials = "email/password invalid"
)

func SignIn(c *fiber.Ctx) error {
	// extract user data from request / JSON
	payload := new(request.SignInRequest)
	if err := c.BodyParser(payload); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, "couldn't parse payload from request: "+err.Error())
	}

	// Validate fields.
	validate := utils.NewValidator()
	if err := validate.Struct(payload); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// get user data from database
	var user entity.User
	if err := database.DB.First(&user, "email = ?", payload.Email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorJSON(c, fiber.StatusUnauthorized, ErrorCredentials)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	// check password validation
	isValid := utils.CheckPasswordHash(payload.Password, user.Password)
	if !isValid {
		return utils.ErrorJSON(c, fiber.StatusUnauthorized, ErrorCredentials)
	}

	// create jwt token
	tokenString, err := utils.CreateToken(&user)
	if err != nil {
		return utils.ErrorJSON(c, fiber.StatusInternalServerError, ErrorCredentials)
	}

	return utils.SuccessJSON(c, fiber.StatusOK, "sign in success", fiber.Map{"token": tokenString})
}
