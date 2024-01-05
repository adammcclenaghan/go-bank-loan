package models

// Requests to the bank API

type CreateLoanRequest struct {
	ID        string `json:"id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
