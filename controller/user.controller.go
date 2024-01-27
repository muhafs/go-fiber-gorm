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
	ErrorNotFound     = "user not found"
	ErrorDuplicate    = "email already taken"
	ErrorParsePayload = "couldn't parse payload"
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
		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "couldn't fetch list user: "+err.Error())
	}

	// return user list data as JSON
	return utils.SuccessJSON(c, fiber.StatusOK, "user list fetched successfully", users)
}

func GetUser(c *fiber.Ctx) error {
	// extract id from req parameter
	id := c.Params("id")

	// get user data from database
	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorJSON(c, fiber.StatusNotFound, ErrorNotFound)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	// return existed user as JSON
	return utils.SuccessJSON(c, fiber.StatusOK, "user fetched successfully", user)
}

func CreateUser(c *fiber.Ctx) error {
	// extract user data from request / JSON
	payload := new(request.User)
	if err := c.BodyParser(payload); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, ErrorParsePayload)
	}

	// run validator
	if vErr := request.Validate[request.User](*payload); vErr != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, vErr)
	}

	// hash password before push to database
	hashedPassword, _ := utils.HashPassword(payload.Password)

	// store image file if exists into the given destination
	nameOrNull, err := utils.StoreFile(c, "avatar", "public/avatars")
	if err != nil {
		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "Couldn't store the file into disk: "+err.Error())
	}

	// assign valid data as entity
	newUser := entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Avatar:   nameOrNull,
	}

	// store user data into database
	if err := database.DB.Create(&newUser).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrorJSON(c, fiber.StatusConflict, ErrorDuplicate)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "couldn't create user: "+err.Error())
	}

	// return stored data as JSON
	return utils.SuccessJSON(c, fiber.StatusCreated, "user created successfully", newUser)
}

func UpdateUser(c *fiber.Ctx) error {
	// extract id from req parameter
	id := c.Params("id")

	// extract user data from request / JSON
	payload := new(request.User)
	if err := c.BodyParser(payload); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, ErrorParsePayload)
	}

	// run validator
	if vErr := request.Validate[request.User](*payload); vErr != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, vErr)
	}

	// get user data from database by ID
	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorJSON(c, fiber.StatusNotFound, ErrorNotFound)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	// hash password before push to database
	hashedPassword, _ := utils.HashPassword(payload.Password)

	// store image file if exists into the given destination
	nameOrNull, err := utils.StoreFile(c, "avatar", "public/avatars")
	if err != nil {
		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "Couldn't store the file into disk: "+err.Error())
	}

	// assign valid data as entity
	updatedUser := entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
		Avatar:   nameOrNull,
	}

	// update the database user with payloads
	if err := database.DB.Model(&user).Updates(&updatedUser).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return utils.ErrorJSON(c, fiber.StatusConflict, ErrorDuplicate)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "couldn't update user: "+err.Error())
	}

	// return the user updated
	return utils.SuccessJSON(c, fiber.StatusCreated, "user data updated successfully", user)
}

func DeleteUser(c *fiber.Ctx) error {
	// get the ID from the req parameter
	id := c.Params("id")

	// get user data from database
	var user entity.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorJSON(c, fiber.StatusNotFound, ErrorNotFound)
		}

		return utils.ErrorJSON(c, fiber.StatusInternalServerError, err.Error())
	}

	// check if user has image
	if user.Avatar.Valid {
		// if has an image then fetch the name of image
		err := utils.RemoveImage("public/avatars", user.Avatar.String)
		if err != nil {
			return utils.ErrorJSON(c, fiber.StatusInternalServerError, "couldn't remove user avatar from disk: "+err.Error())
		}
	}

	// delete user from database
	result := database.DB.Delete(&user)

	// if there is no record matches, it means there is no user exists
	if result.RowsAffected == 0 {
		return utils.ErrorJSON(c, fiber.StatusNotFound, ErrorNotFound)
	} else if result.Error != nil {
		return utils.ErrorJSON(c, fiber.StatusInternalServerError, "couldn't delete user: "+result.Error.Error())
	}

	// return success message
	return utils.SuccessJSON(c, fiber.StatusOK, "user deleted successfully", true)
}
