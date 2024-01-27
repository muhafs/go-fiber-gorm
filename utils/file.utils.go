package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func StoreFile(c *fiber.Ctx, fieldKey string, dirPath string) (sql.NullString, error) {
	// get file image from avatar form
	file, _ := c.FormFile(fieldKey)

	// store image if exists
	if file != nil {
		// check if directory exists, if not create it first.
		err := CheckDirectory(dirPath)
		if err != nil {
			return sql.NullString{}, err
		}

		// generate new file name
		fileName := GenerateUniqueName(file)

		// store image to destination path
		destination := fmt.Sprintf("./%s/%s", dirPath, fileName)
		if err := c.SaveFile(file, destination); err != nil {
			return sql.NullString{}, err
		}

		// return the string name
		return sql.NullString{Valid: true, String: fileName}, nil
	}

	// otherwise return null string
	return sql.NullString{}, nil
}

func StoreFiles(c *fiber.Ctx, fieldKey string, dirPath string) ([]sql.NullString, error) {
	// Parse the multipart form
	form, _ := c.MultipartForm()

	// Get all files from field key
	files := form.File[fieldKey]

	// Loop through files:
	var storedFiles []sql.NullString
	for i, file := range files {
		// generate new file name
		fileName := GenerateUniqueName(file)

		// Save the files to disk:
		destination := fmt.Sprintf("./%s/%v_%s", dirPath, i+1, fileName)
		if err := c.SaveFile(file, destination); err != nil {
			return nil, err
		}

		// append file name into temp store
		storedFiles = append(storedFiles, sql.NullString{Valid: true, String: fileName})
	}

	// return all file names
	return storedFiles, nil
}

func CheckDirectory(dirPath string) error {
	// create directory if not exists
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	// if exists just skip it
	return nil
}

func GenerateUniqueName(file *multipart.FileHeader) string {
	// generate new uuid
	// eg. "a11-b22-c33-d44"
	uniqueId := uuid.New()

	// remove "-" from unique id
	// eg. "a11b22c33d44"
	convertedID := strings.Replace(uniqueId.String(), "-", "", -1)

	// generate unix from current time
	// eg. "1654041600"
	currentUnix := time.Now().Unix()

	// concat unix and unique id by "_"
	// eg. "1654041600_a11b22c33d44"
	filename := fmt.Sprintf("%v_%s", currentUnix, convertedID)

	// extract image extension from original file filename
	// eg. "png"
	splitedName := strings.Split(file.Filename, ".")
	fileExt := splitedName[len(splitedName)-1]

	// generate image from filename and extension
	// eg. "a11b22c33d44.png"
	return fmt.Sprintf("%s.%s", filename, fileExt)
}

func RemoveImage(dirPath, fileImage string) error {
	// set image path to remove
	imagePath := fmt.Sprintf("./%s/%s", dirPath, fileImage)

	// remove image from disk
	if err := os.Remove(imagePath); err != nil {
		return err
	}

	// return if success
	return nil
}
