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

func TestNotesController_Integration(t *testing.T) {
	db := SetupTestDB()
	r := SetupTestRouter(db)

	// 1. Register and Login to get token
	body := map[string]string{
		"email":    "notes@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["token"]

	t.Run("Create Note", func(t *testing.T) {
		noteBody := map[string]string{"title": "Integration Note"}
		jsonNote, _ := json.Marshal(noteBody)
		req, _ := http.NewRequest("POST", "/api/v1/notes/", bytes.NewBuffer(jsonNote))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response map[string]models.Note
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Integration Note", response["data"].Title)
	})

	t.Run("Get Notes", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/notes/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string][]models.Note
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response["data"], 1)
	})
}
