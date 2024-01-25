package services

import (
	"api/models"
	"api/repository"
	"api/store"
	"api/utils"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrFailedGetProspectUser     = errors.New("failed to get prospect user")
	ErrFailedCheckProspectSwiped = errors.New("failed to check if prospect swiped back")
	ErrFailedCreateSwipe         = errors.New("failed to create swipe")
	ErrFailedCreateMatch         = errors.New("failed to create match")
	ErrFailedGenerateID          = errors.New("failed to generate ID")
)

type SwipeService struct {
	eventStore      store.EventStore
	logger          *logrus.Logger
	swipeRepository repository.SwipesRepository
	matchRepository repository.MatchRepository
	userRepository  repository.UserRepository
}

func NewSwipeService(eventStore store.EventStore, logger *logrus.Logger, swipeRepository repository.SwipesRepository,
	matchRepository repository.MatchRepository,
	userRepository repository.UserRepository,
) *SwipeService {
	return &SwipeService{
		eventStore:      eventStore,
		logger:          logger,
		swipeRepository: swipeRepository,
		matchRepository: matchRepository,
		userRepository:  userRepository,
	}
}

// Swipe swipes a user through a prospect profile for a possible match.
func (s *SwipeService) Swipe(ctx context.Context, userID string, payload models.SwipePayload) (*models.SwipeResponse, error) {
	_, err := s.userRepository.GetUserById(ctx, payload.ProspectID)
	if err != nil {
		s.logger.WithError(err).Error(ErrFailedGetProspectUser)
		return nil, ErrFailedGetProspectUser
	}

	checkIfImProspectSwipe, err := s.swipeRepository.GetSwipeByUserAndProspect(ctx, payload.ProspectID, userID)
	if err != nil {
		s.logger.WithError(err).Error(ErrFailedCheckProspectSwiped)
		return nil, ErrFailedCheckProspectSwiped
	}

	swipe := &models.Swipe{
		ID:         utils.GenerateId(),
		UserID:     userID,
		ProspectID: payload.ProspectID,
		Interested: payload.Interested,
	}

	swipe, err = s.swipeRepository.CreateSwipe(ctx, swipe)
	if err != nil {
		s.logger.WithError(err).Error(ErrFailedCreateSwipe)
		return nil, ErrFailedCreateSwipe
	}

	var matchUser *models.Match

	if checkIfImProspectSwipe != nil && checkIfImProspectSwipe.Interested && swipe.Interested {
		match := &models.Match{
			ID:       utils.GenerateId(),
			Profiles: []string{userID, payload.ProspectID},
			Matched:  true,
		}

		matchUser, err = s.matchRepository.CreateMatch(ctx, match)
		if err != nil {
			s.logger.WithError(err).Error(ErrFailedCreateMatch)
			return nil, ErrFailedCreateMatch
		}
	}

	if matchUser == nil {
		return &models.SwipeResponse{
			Matched: false,
		}, nil
	}

	return &models.SwipeResponse{
		Matched: true,
		MatchID: matchUser.ID,
	}, nil
}
