package controllers

import (
	"backend/models"
	"backend/services"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadUserBackups(c *gin.Context) {
	claims := c.MustGet("claims").(*models.AccessTokenClaims)

	uploadDir := "../uploads/"
	dbDir := uploadDir + "db/"
	imgDir := uploadDir + "images/"

	os.MkdirAll(dbDir, os.ModePerm)
	os.MkdirAll(imgDir, os.ModePerm)

	// ---- Receive DB file ----
	dbFile, dbHeader, err := c.Request.FormFile("db_file")
	if err != nil {
		c.JSON(400, gin.H{"message": "db_file not found"})
		return
	}
	defer dbFile.Close()

	// ---- Receive Image file ----
	imgFile, imgHeader, err := c.Request.FormFile("image_file")
	if err != nil {
		c.JSON(400, gin.H{"message": "image_file not found"})
		return
	}
	defer imgFile.Close()

	// Safe filenames
	dbFileName := fmt.Sprintf("%d_%s", claims.UserID, filepath.Base(dbHeader.Filename))
	imgFileName := fmt.Sprintf("%d_%s", claims.UserID, filepath.Base(imgHeader.Filename))

	dbPath := dbDir + dbFileName
	imgPath := imgDir + imgFileName

	// ---- Check existing backup ----
	var backup models.DataBackup
	result := services.DB.Where("user_id = ?", claims.UserID).First(&backup)

	// ---- Save DB file ----
	dbDst, err := os.Create(dbPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot save db file"})
		return
	}
	_, err = io.Copy(dbDst, dbFile)
	dbDst.Close()
	if err != nil {
		c.JSON(500, gin.H{"error": "db upload failed"})
		return
	}

	// ---- Save Image file ----
	imgDst, err := os.Create(imgPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot save image file"})
		return
	}
	_, err = io.Copy(imgDst, imgFile)
	imgDst.Close()
	if err != nil {
		c.JSON(500, gin.H{"error": "image upload failed"})
		return
	}

	// ---- Delete old files AFTER successful save ----
	if result.Error == nil {
		if backup.UserDBFiles != "" {
			_ = os.Remove(dbDir + backup.UserDBFiles)
		}
		if backup.UserImages != "" {
			_ = os.Remove(imgDir + backup.UserImages)
		}
	}

	// ---- Insert or Update DB ----
	if result.Error == nil {
		services.DB.Model(&models.DataBackup{}).
			Where("user_id = ?", claims.UserID).
			Updates(map[string]interface{}{
				"user_dbfiles": dbFileName,
				"user_images":  imgFileName,
			})
	} else {
		newBackup := models.DataBackup{
			UserID:      claims.UserID,
			UserDBFiles: dbFileName,
			UserImages:  imgFileName,
		}
		services.DB.Create(&newBackup)
	}

	c.JSON(200, gin.H{
		"message":  "backup uploaded",
		"db_file":  dbFileName,
		"img_file": imgFileName,
	})
}

func RestoreDatabase(c *gin.Context) {
	claims := c.MustGet("claims").(*models.AccessTokenClaims)

	var backup models.DataBackup
	err := services.DB.Where("user_id = ?", claims.UserID).First(&backup).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "no backup found"})
		return
	}

	filePath := "../uploads/db/" + backup.UserDBFiles

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"message": "database file missing"})
		return
	}

	c.File(filePath) // Streams directly (supports large files)
}

func RestoreImages(c *gin.Context) {
	claims := c.MustGet("claims").(*models.AccessTokenClaims)

	var backup models.DataBackup
	err := services.DB.Where("user_id = ?", claims.UserID).First(&backup).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "no backup found"})
		return
	}

	filePath := "../uploads/images/" + backup.UserImages

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"message": "image file missing"})
		return
	}

	c.File(filePath)
}
