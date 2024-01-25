package setup

import (
	"api/config"
	"api/middlewares"
	"api/repository"
	"api/repository/mongodb"
	"api/services"
	"api/store"
	"fmt"
	"github.com/sirupsen/logrus"
)

type ServiceDependencies struct {
	EventStore   store.EventStore
	Logger       *logrus.Logger
	UserService  *services.UserService
	SwipeService *services.SwipeService
	MatchService *services.MatchService
	Middlewares  *middlewares.SystemMiddleware
}

// ServiceInitializer is an interface for initializing services.
type ServiceInitializer interface {
	Init() (*ServiceDependencies, error)
}

// MongoDBInitializer initializes MongoDB-related services.
type MongoDBInitializer struct {
	Secrets config.Secrets
	Logger  *logrus.Logger
}

func (m MongoDBInitializer) Init() (*ServiceDependencies, error) {
	var (
		eventStore      store.EventStore
		userRepository  repository.UserRepository
		swipeRepository repository.SwipesRepository
		matchRepository repository.MatchRepository
	)

	m.Logger.Info("Using MongoDB as the database")
	mongoStore, err := mongodb.NewMongoConnection(m.Secrets.DatabaseUrl, m.Secrets.DatabaseName)
	if err != nil {
		return nil, fmt.Errorf("error opening MongoDB database: %w", err)
	}
	userRepository = mongodb.NewUserRepo(mongoStore)
	swipeRepository = mongodb.NewSwipeRepo(mongoStore)
	matchRepository = mongodb.NewMatchRepo(mongoStore)

	eventStore = store.NewEventStore(m.Logger)

	userService := services.NewUserService(eventStore, userRepository, m.Logger, m.Secrets.JwtSecret)

	return &ServiceDependencies{
		EventStore:   eventStore,
		Logger:       m.Logger,
		UserService:  services.NewUserService(eventStore, userRepository, m.Logger, m.Secrets.JwtSecret),
		SwipeService: services.NewSwipeService(eventStore, m.Logger, swipeRepository, matchRepository, userRepository),
		MatchService: services.NewMatchService(eventStore, m.Logger, matchRepository),
		Middlewares:  middlewares.NewSystemMiddleware(userService, m.Logger),
	}, nil
}

func ConfigureServiceDependencies(initializer ServiceInitializer) (*ServiceDependencies, error) {
	return initializer.Init()
}
