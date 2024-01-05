package shared_models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusPendingIsValid(t *testing.T) {
	statusStr := "pending"
	status := Status(statusStr)

	assert.True(t, status.IsValid())
}

func TestStatusCompletedIsValid(t *testing.T) {
	statusStr := "completed"
	status := Status(statusStr)

	assert.True(t, status.IsValid())
}

func TestStatusRejectedIsValid(t *testing.T) {
	statusStr := "rejected"
	status := Status(statusStr)

	assert.True(t, status.IsValid())
}

func TestStatusUnexpectedIsInvalid(t *testing.T) {
	statusStr := "some unexpected status"
	status := Status(statusStr)

	assert.False(t, status.IsValid())
}