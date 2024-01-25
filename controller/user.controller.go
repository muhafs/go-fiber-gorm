package controller

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-gorm/database"
	"github.com/muhafs/go-fiber-gorm/model/entity"
	"github.com/muhafs/go-fiber-gorm/model/request"
	"github.com/muhafs/go-fiber-gorm/utils"
	"gorm.io/gorm"
)

const (
	ErrorNotFound  = "User not found"
	ErrorDuplicate = "email already taken"
)

func GetListUser(c *fiber.Ctx) error {
	// extract the page and limit query
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	// convert string into int, so we can calculate the offset page
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	// fetch user list from database
	var users []entity.User
	if err := database.DB.Limit(intLimit).Offset(offset).Order("name asc").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "couldn't fetch list user: " + err.Error(),
		})
	}

	// return user list data as JSON
	return c.JSON(fiber.Map{
		"success": true,
		"message": "user list fetched successfully",
		"data":    users,
	})
}

func GetUser(c *fiber.Ctx) error {
	// extract id from req parameter
	id := c.Params("id")

	// get user data from database
	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": ErrorNotFound,
			})
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// return existed user as JSON
	return c.JSON(fiber.Map{
		"success": true,
		"message": "user fetched successfully",
		"data":    user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	// extract user data from request / JSON
	payload := new(request.User)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "couldn't parse payload from request: " + err.Error(),
		})
	}

	// run validator
	if vErr := payload.Validate(); vErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": vErr,
		})
	}

	pass, _ := utils.HashPassword(payload.Password)
	// assign valid data as entity
	newUser := entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: pass,
	}

	// store user data into database
	if err := database.DB.Create(&newUser).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": ErrorDuplicate,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "couldn't create user: " + err.Error(),
		})
	}

	// return stored data as JSON
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "user created successfully",
		"data":    newUser,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	// extract id from req parameter
	id := c.Params("id")

	// extract user data from request / JSON
	payload := new(request.User)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "couldn't parse payload from request: " + err.Error(),
		})
	}

	// run validator
	if vErr := payload.Validate(); vErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": vErr,
		})
	}

	// get user data from database by ID
	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": ErrorNotFound,
			})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// assign valid data as entity
	pass, _ := utils.HashPassword(payload.Password)
	updatedUser := entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: pass,
	}

	// update the database user with payloads
	if err := database.DB.Model(&user).Updates(&updatedUser).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": ErrorDuplicate,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "couldn't update user: " + err.Error(),
		})
	}

	// return the user updated
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "user data updated successfully",
		"data":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	// get the ID from the req parameter
	id := c.Params("id")

	// delete user from database with the given ID
	result := database.DB.Delete(&entity.User{}, id)

	// if there is no record matches, it means there is no user exists
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": ErrorNotFound,
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"message": "couldn't delete user: " + result.Error.Error(),
		})
	}

	// return success message
	return c.JSON(fiber.Map{
		"success": true,
		"message": "user deleted successfully",
	})
}
