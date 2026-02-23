package controllers

import (
	"backend/models"
	"backend/services"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EmailVerification(c *gin.Context) {
	// variable to store the request body.
	var requestBody struct {
		UserEmail string `json:"user_email"`
	}

	// Binded the requestBody variable to the request.
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON Body.",
		})
		return
	}

	if requestBody.UserEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	var emailCount int64

	services.DB.Model(&models.Users{}).Where("user_email = ?", requestBody.UserEmail).Count(&emailCount)
	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email already registered"})
		return
	}

	otp := services.GenerateOTP()

	services.SaveOTP(requestBody.UserEmail, otp)
	go func() {
		if err := services.SendOtp(requestBody.UserEmail, otp); err != nil {
			log.Println("[ASYNC EMAIL ERROR]", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP has been sent.",
	})
}

func OtpVerification(c *gin.Context) {
	var requestBody struct {
		UserEmail string `json:"user_email"`
		Otp       string `json:"otp"`
		Purpose   string `json:"purpose"` // NEW
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON Body"})
		return
	}

	if requestBody.UserEmail == "" || requestBody.Otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email, OTP are required"})
		return
	}

	if requestBody.Purpose != "email_verification" && requestBody.Purpose != "password_reset" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid purpose"})
		return
	}

	if !services.VerifyOTP(requestBody.UserEmail, requestBody.Otp) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired OTP"})
		return
	}

	services.DeleteOTP(requestBody.UserEmail)

	token, err := services.GenerateEmailVerificationToken(requestBody.UserEmail, requestBody.Purpose)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":            "OTP verified",
		"verification_token": token,
	})
}

func UserSignup(c *gin.Context) {
	var requestBody struct {
		UserEmail         string `json:"user_email"`
		UserFullName      string `json:"user_fullname"`
		Password          string `json:"password"`
		Dob               string `json:"user_dob"`
		Gender            string `json:"gender"`
		VerificationToken string `json:"verification_token"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	claims, err := services.ValidateToken(requestBody.VerificationToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired verification token"})
		return
	}

	if claims["user_email"] != requestBody.UserEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email mismatch"})
		return
	}

	hashedPassword, err := services.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	dob, err := time.ParseInLocation("2006-01-02", requestBody.Dob, time.UTC)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid DOB format. Use YYYY-MM-DD"})
		return
	}

	user := models.Users{
		UserEmail:    requestBody.UserEmail,
		UserFullName: requestBody.UserFullName,

		PasswordHash: hashedPassword,
		Gender:       requestBody.Gender,

		Dob: dob,
	}

	if err := services.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User creation failed",
			"error":   err.Error(),
		})
		return
	}

	jwtToken, err := services.GenerateAuthenticationToken(uint(user.UserID), user.UserEmail)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "JWT token generation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User created successfully",
		"access_token": jwtToken,
	})
}

func UserSignin(c *gin.Context) {
	var requestBody struct {
		UserEmail string `json:"user_email"`
		Password  string `json:"password"`
	}

	var user models.Users

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON",
		})
		return
	}

	if requestBody.UserEmail == "" || requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	if err := services.DB.Where("user_email = ?", requestBody.UserEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	if services.ComparePasswordAndHashedPassword(user.PasswordHash, requestBody.Password) {
		jwtToken, err := services.GenerateAuthenticationToken(uint(user.UserID), user.UserEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":      "Login successful",
			"access_token": jwtToken,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Incorrect Password",
	})
}

func ForgotPassword(c *gin.Context) {
	var requestBody struct {
		UserEmail string `json:"user_email"`
	}

	var expectedUser models.Users

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON Body.",
		})
		return
	}

	if requestBody.UserEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	if err := services.DB.Where("user_email = ?", requestBody.UserEmail).First(&expectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email is not registered"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	otp := services.GenerateOTP()

	services.SaveOTP(requestBody.UserEmail, otp)

	go func() {
		if err := services.SendOtp(requestBody.UserEmail, otp); err != nil {
			log.Println("[ASYNC EMAIL ERROR]", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP has been sent",
	})
}

func ResetPassword(c *gin.Context) {
	var requestBody struct {
		VerificationToken string `json:"verification_token"`
		NewPassword       string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON Body.",
		})
		return
	}

	claims, err := services.ValidateToken(requestBody.VerificationToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"messagesss": err.Error()})
		return
	}

	if claims["purpose"] != "password_reset" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token type"})
		return
	}

	var user models.Users
	email := claims["user_email"].(string)

	if err := services.DB.Where("user_email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	hashedPassword, err := services.HashPassword(requestBody.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error hashing password",
		})
		return
	}

	err = services.DB.Model(&user).Update("password_hash", hashedPassword).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func UpdatePassword(c *gin.Context) {
	claims := c.MustGet("claims").(*models.AccessTokenClaims)

	verificationToken, err := services.GenerateEmailVerificationToken(
		claims.UserEmail,
		"password_reset",
	)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message":            "password reset token generated",
		"verification_token": verificationToken,
	})
}

func UpdateUser(c *gin.Context) {
	var requestBody struct {
		UserFullName string `json:"user_fullname"`
		Gender       string `json:"gender"`
		Dob          string `json:"user_dob"`
	}

	// 1. Bind request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// 2. Extract claims (trusted identity)
	claims := c.MustGet("claims").(*models.AccessTokenClaims)
	userEmail := claims.UserEmail

	// 3. Build update map safely
	updates := map[string]interface{}{}

	if requestBody.UserFullName != "" {
		updates["user_fullname"] = requestBody.UserFullName
	}

	if requestBody.Gender != "" {
		updates["gender"] = requestBody.Gender
	}

	if requestBody.Dob != "" {
		updates["user_dob"] = requestBody.Dob
	}

	if len(updates) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "no fields to update"})
		return
	}

	// 4. Update database
	result := services.DB.Model(&models.Users{}).
		Where("user_email = ?", userEmail).
		Updates(updates)

	if result.Error != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "failed to update user"})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"message": "profile updated successfully",
	})
}

func GetUserProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*models.AccessTokenClaims)

	type UserProfileResponse struct {
		UserFullname string  `json:"user_fullname"`
		UserDob      *string `json:"user_dob"`
		Gender       string  `json:"gender"`
	}

	// Temporary struct to read from DB
	var dbResult struct {
		UserFullname string
		UserDob      *time.Time
		Gender       string
	}

	err := services.DB.
		Model(&models.Users{}).
		Select("user_fullname, user_dob, gender").
		Where("user_email = ?", claims.UserEmail).
		First(&dbResult).Error

	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	var dobStr *string
	if dbResult.UserDob != nil {
		formatted := dbResult.UserDob.Format("2006-01-02") // YYYY-MM-DD
		dobStr = &formatted
	}

	profile := UserProfileResponse{
		UserFullname: dbResult.UserFullname,
		UserDob:      dobStr,
		Gender:       dbResult.Gender,
	}

	c.JSON(200, profile)
}
