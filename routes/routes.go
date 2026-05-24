package routes

import (
	"gin-test/controllers"
	"gin-test/middleware"
	"gin-test/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	// Static files
	r.Static("/uploads", "./uploads")

	// Services
	authService := services.NewAuthService(db)
	notesService, _ := services.NewNotesService(db)

	// Controllers
	authController := controllers.NewAuthController(authService)
	notesController := controllers.NewNoteController(notesService)
	shareController := controllers.NewShareController(notesService)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		// Public Auth Routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected Notes Routes
		notes := v1.Group("/notes")
		notes.Use(middleware.AuthMiddleware())
		{
			notes.GET("/", notesController.GetNotes)
			notes.POST("/", notesController.CreateNote)
			notes.GET("/:id", notesController.GetNoteByID)
			notes.PUT("/:id", notesController.UpdateNote)
			notes.DELETE("/:id", notesController.DeleteNote)
			notes.POST("/:id/attachments", notesController.UploadAttachment)
			notes.POST("/:id/share", shareController.ShareNote)
			notes.DELETE("/:id/share/:uid", shareController.RevokeShare)
			notes.GET("/:id/shares", shareController.ListShares)
		}
	}

	return r
}
