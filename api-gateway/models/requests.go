//Package models provides models used by controllers
package models

// CreateApplicationRequest represents an API request to create a new loan application
type CreateApplicationRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
