package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"drgo/internal/middleware"
	medical_router "drgo/internal/router/patient-panel"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to DrGo API",
		})
	})

	api.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	publicRoutes := api.Group("/public")
	medical_router.SetupPatientPanelRoutes(publicRoutes, db)

	return r
}
