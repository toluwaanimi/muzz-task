package repository

import (
	"api/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) InsertUsers(ctx context.Context, users []*models.User) error {
	args := m.Called(ctx, users)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, payload *models.User) (*models.User, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, payload *models.User) (*models.User, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) Discover(ctx context.Context, filter models.UserFilter, user models.User) ([]*models.User, error) {
	args := m.Called(ctx, filter, user)
	return args.Get(0).([]*models.User), args.Error(1)
}

type MockSwipeRepository struct {
	mock.Mock
}

func (m *MockSwipeRepository) GetSwipeByUserAndProspect(ctx context.Context, prospectID, userID string) (*models.Swipe, error) {
	args := m.Called(ctx, prospectID, userID)
	return args.Get(0).(*models.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) CreateSwipe(ctx context.Context, swipe *models.Swipe) (*models.Swipe, error) {
	args := m.Called(ctx, swipe)
	return args.Get(0).(*models.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) GetSwipeById(ctx context.Context, id string) (*models.Swipe, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) UpdateSwipe(ctx context.Context, payload *models.Swipe) (*models.Swipe, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.Swipe), args.Error(1)
}

func (m *MockSwipeRepository) DeleteSwipe(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) CreateMatch(ctx context.Context, match *models.Match) (*models.Match, error) {
	args := m.Called(ctx, match)
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchRepository) GetMatchById(ctx context.Context, id string) (*models.Match, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchRepository) GetMatchesFiltered(ctx context.Context, filter models.MatchFilter) (*models.MatchedUserInfo, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.MatchedUserInfo), args.Error(1)
}

func (m *MockMatchRepository) UpdateMatch(ctx context.Context, payload *models.Match) (*models.Match, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchRepository) DeleteMatch(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
