package services

import (
	"api/models"
	"api/repository"
	"api/store"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSwipeService_Swipe(t *testing.T) {
	// Mock repositories
	userRepo := new(repository.MockUserRepository)
	swipeRepo := new(repository.MockSwipeRepository)
	matchRepo := new(repository.MockMatchRepository)

	// Setup service
	eventStore := store.NewEventStore(logrus.New())
	logger := logrus.New()

	swipeService := NewSwipeService(eventStore, logger, swipeRepo, matchRepo, userRepo)

	// Test case data
	userID := "user123"
	prospectID := "prospect456"
	payload := models.SwipePayload{
		ProspectID: prospectID,
		Interested: true,
	}

	// Mock repository expectations
	userRepo.On("GetUserById", mock.Anything, prospectID).Return(&models.User{}, nil)
	swipeRepo.On("GetSwipeByUserAndProspect", mock.Anything, prospectID, userID).Return(&models.Swipe{}, nil)
	swipeRepo.On("CreateSwipe", mock.Anything, mock.AnythingOfType("*models.Swipe")).Return(&models.Swipe{}, nil)
	matchRepo.On("CreateMatch", mock.Anything, mock.AnythingOfType("*models.Match")).Return(&models.Match{}, nil).Maybe() // Adjusted to expect exactly one call

	// Perform the swipe
	response, err := swipeService.Swipe(context.Background(), userID, payload)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.Matched)

	// Verify repository method calls
	userRepo.AssertExpectations(t)
	swipeRepo.AssertExpectations(t)
	matchRepo.AssertExpectations(t)
}
