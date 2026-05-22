package errors

import "net/http"

type AppError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

func (e *AppError) Error() string { return e.Message }

func NotFound(msg string) *AppError      { return &AppError{http.StatusNotFound, msg} }
func Forbidden(msg string) *AppError     { return &AppError{http.StatusForbidden, msg} }
func BadRequest(msg string) *AppError    { return &AppError{http.StatusBadRequest, msg} }
func Conflict(msg string) *AppError      { return &AppError{http.StatusConflict, msg} }
func Internal(msg string) *AppError      { return &AppError{http.StatusInternalServerError, msg} }
func Unauthorized(msg string) *AppError  { return &AppError{http.StatusUnauthorized, msg} }
