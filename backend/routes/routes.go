package routes

import (
	"coworking-backend/handlers"
	"coworking-backend/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/health", handlers.Health)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.RequireAuth(), handlers.Me)
		}

		spaces := api.Group("/spaces")
		{
			spaces.GET("", middleware.RequireAuth(), handlers.ListSpaces)
			spaces.GET("/random", middleware.RequireAuth(), handlers.RandomAvailableSpace)
			spaces.POST("", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.CreateSpace)
			spaces.DELETE("/:id", middleware.RequireAuth(), middleware.RequireAdmin(), handlers.DeleteSpace)
		}

		bookings := api.Group("/bookings")
		bookings.Use(middleware.RequireAuth())
		{
			bookings.POST("", handlers.CreateBooking)
			bookings.GET("/me", handlers.ListMyBookings)
			bookings.GET("/streak", handlers.Streak)
			bookings.PATCH("/:id/cancel", handlers.CancelBooking)
		}
	}
}
