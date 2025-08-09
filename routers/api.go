package routers

import (
	"github.com/gin-gonic/gin"
	"os"
	"spe-trx-gateway/controllers"
	"spe-trx-gateway/middlewares"
)

func WebRoute(router *gin.Engine, s *controllers.Server) {
	internal := router.Group("/api/v1/internal")
	{
		internal.POST("/hash-notification")
		internal.POST("/hash-inquiry")
	}

	register := router.Group("/api/v1/register")
	{
		register.POST("/get-auth", controllers.RegisterController)
	}

	payment := router.Group("/api/v1/payment")
	payment.Use(middlewares.SignatureMiddleware())
	{
		payment.POST("/transaction-notification", s.NotificationController)
		payment.POST("/check-status", s.InquiryController)
	}
}

// Route is the driver for the webapi service.
func Route(s *controllers.Server) (*gin.Engine, string) {
	appRelease := os.Getenv("APP_MODE")
	if appRelease == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger(), CORS())

	WebRoute(router, s)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "9191"
	}

	return router, port
}

// CORS Cross Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, "+
			"Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
