package services

import (
	"api/repository"
	"api/store"
	"github.com/sirupsen/logrus"
)

type MatchService struct {
	eventStore      store.EventStore
	logger          *logrus.Logger
	matchRepository repository.MatchRepository
}

func NewMatchService(eventStore store.EventStore, logger *logrus.Logger, matchRepository repository.MatchRepository) *MatchService {
	return &MatchService{
		eventStore:      eventStore,
		logger:          logger,
		matchRepository: matchRepository,
	}
}
