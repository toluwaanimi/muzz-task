package repository

import (
	"api/models"
	"context"
	"errors"
	"github.com/olivere/elastic/v7"
)

var ErrDuplicateFound = errors.New("error: duplicate found")

type UserRepository interface {
	InsertUsers(ctx context.Context, users []*models.User) error
	CreateUser(ctx context.Context, payload *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, payload *models.User) (*models.User, error)
	GetUserCount(ctx context.Context) (int, error)
	Discover(ctx context.Context, filter models.UserFilter, user models.User) ([]*models.User, error)
}

type SwipesRepository interface {
	CreateSwipe(ctx context.Context, payload *models.Swipe) (*models.Swipe, error)
	GetSwipeById(ctx context.Context, id string) (*models.Swipe, error)
	GetSwipeByUserAndProspect(ctx context.Context, userID, prospectID string) (*models.Swipe, error)
	UpdateSwipe(ctx context.Context, payload *models.Swipe) (*models.Swipe, error)
	DeleteSwipe(ctx context.Context, id string) error
}

type MatchRepository interface {
	CreateMatch(ctx context.Context, payload *models.Match) (*models.Match, error)
	GetMatchById(ctx context.Context, id string) (*models.Match, error)
	GetMatchesFiltered(ctx context.Context, filter models.MatchFilter) (*models.MatchedUserInfo, error)
	UpdateMatch(ctx context.Context, payload *models.Match) (*models.Match, error)
	DeleteMatch(ctx context.Context, id string) error
}

type ElasticsearchRepository interface {
	Search(ctx context.Context, indexName string, query *elastic.BoolQuery, resultType interface{}) ([]interface{}, error)
	Index(ctx context.Context, index string, id string, document interface{}) error
}
