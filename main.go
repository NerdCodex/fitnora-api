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

	userRoute.POST("/update", controllers.UpdateUser)
	userRoute.POST("/updatepassword", controllers.UpdatePassword)
	userRoute.GET("/profile", controllers.GetUserProfile)

	router.Run("10.55.230.72:8080")
}
