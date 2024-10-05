package main

import (
	controllers "gin-test/Controllers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// router.POST("/me", func(c *gin.Context) {

	// 	type ReqBody struct {
	// 		Email string `json:"email" binding:"required"`
	// 		Name  string `json:"name"`
	// 	}

	// 	var reqBody ReqBody

	// 	if err := c.BindJSON(&reqBody); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"Email": reqBody.Email,
	// 		"Name":  reqBody.Name,
	// 	})
	// })

	notesController := &controllers.NoteControllers{}
	notesController.InitNotesControllersRoutes(router)

	router.Run(":8000")

}
