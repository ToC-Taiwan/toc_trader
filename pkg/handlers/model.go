// Package handlers handlers
package handlers

// ErrorResponse ErrorResponse
type ErrorResponse struct {
	Response   string      `json:"response"`
	Attachment interface{} `json:"attachment"`
}
