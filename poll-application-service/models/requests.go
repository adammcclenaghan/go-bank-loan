//Package models represents models for HTTP data
package models

type PollLoanRequest struct {
	ApplicationID string `json:"application_id" binding:"required"`
}
