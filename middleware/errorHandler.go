package middleware

import (
	"net/http"

	apperrors "gin-test/internal/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*apperrors.AppError); ok {
			c.JSON(appErr.Status, appErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
