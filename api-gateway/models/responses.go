//Package models provides models used by controllers
package models

import sharedmodels "service-shared/shared-models"

// CreateApplicationResponse represents an API response to a CreateApplicationRequest
type CreateApplicationResponse struct {
	ApplicationID string              `json:"application_id" binding:"required""`
	Status        sharedmodels.Status `json:"status" binding:"required"`
	FirstName     string              `json:"first_name" binding:"required"`
	LastName      string              `json:"last_name" binding:"required"`
}

// ClientApplicationView represents the information we provide to clients of this API for a loan application
type ClientApplicationView struct {
	ApplicationID string              `json:"application_id" binding:"required" bson:"_id"`
	Status        sharedmodels.Status `json:"status" binding:"required, validstatus" bson:"status"`
	FirstName     string              `json:"first_name" binding:"required"`
	LastName      string              `json:"last_name" binding:"required"`
}

// GetAppsWithStatusResponse provides the client with a view of all applications with a given status
type GetAppsWithStatusResponse struct {
	ApplicationsWithStatus []ClientApplicationView `json:"applications" binding:"required"`
}
