package integration

import (
	"bytes"
	"encoding/json"
	"gin-test/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthController_Integration(t *testing.T) {
	db := SetupTestDB()
	r := SetupTestRouter(db)

	t.Run("Register Success", func(t *testing.T) {
		body := map[string]string{
			"email":    "reg@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response map[string]models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "reg@example.com", response["data"].Email)
	})

	t.Run("Login Success", func(t *testing.T) {
		body := map[string]string{
			"email":    "reg@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
	})
}
