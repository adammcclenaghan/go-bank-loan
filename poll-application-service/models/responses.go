//Package models represents models for HTTP bodies.
package models

//PollLoanResponse represents a response from the jobs endpoint of the bank API
type PollLoanResponse struct {
	ApplicationID string `json:"application_id" binding:"required"`
	Status        string `json:"status" binding:"required"`
}
