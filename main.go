package main

import (
	"backend/controllers"
	"backend/middleware"
	"backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	services.ConnectToDB()

	router := gin.Default()

	userRoute := router.Group("/user")
	userRoute.Use(middleware.ValidateJWT)

	router.POST("email", controllers.EmailVerification)
	router.POST("otp", controllers.OtpVerification)
	router.POST("signup", controllers.UserSignup)
	router.POST("signin", controllers.UserSignin)
	router.POST("forgotpassword", controllers.ForgotPassword)
	router.POST("resetpassword", controllers.ResetPassword)
	router.POST("example", controllers.AnalyzeFood)

	userRoute.POST("/update", controllers.UpdateUser)
	userRoute.GET("/updatepassword", controllers.UpdatePassword)
	userRoute.GET("/profile", controllers.GetUserProfile)

	// Data Backup routes
	userRoute.POST("/backup/upload", controllers.UploadUserBackups)
	userRoute.GET("/backup/database", controllers.RestoreDatabase)
	userRoute.GET("/backup/images", controllers.RestoreImages)

	// Food Analysis routes
	userRoute.POST("/food/analyze", controllers.AnalyzeFood)

	router.Run("localhost:8080")
}
