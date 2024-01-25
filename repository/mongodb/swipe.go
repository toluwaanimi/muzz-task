package mongodb

import (
	"api/constants"
	"api/models"
	"api/repository"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type swipeRepository struct {
	mongo          *MongoStore
	collection     string
	userCollection string
}

// GetSwipeByUserAndProspect returns a swipe by the given user and prospect IDs.
func (s swipeRepository) GetSwipeByUserAndProspect(ctx context.Context, userID, prospectID string) (*models.Swipe, error) {
	var swipe models.Swipe
	err := s.mongo.coll(s.collection).FindOne(ctx, bson.M{"user_id": userID, "prospect_id": prospectID}).Decode(&swipe)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &swipe, nil
}

// CreateSwipe creates a new swipe in the database.
func (s swipeRepository) CreateSwipe(ctx context.Context, payload *models.Swipe) (*models.Swipe, error) {
	existingSwipe, err := s.GetSwipeByUserAndProspect(ctx, payload.UserID, payload.ProspectID)
	if err != nil {
		return nil, err
	}

	if existingSwipe != nil {
		return nil, repository.ErrDuplicateFound
	}
	_, err = s.mongo.coll(s.collection).InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// GetSwipeById returns a swipe by their ID.
func (s swipeRepository) GetSwipeById(ctx context.Context, id string) (*models.Swipe, error) {
	var swipe models.Swipe
	if err := s.mongo.coll(s.collection).FindOne(ctx, bson.M{"id": id}).Decode(&swipe); err != nil {
		return nil, err
	}
	return &swipe, nil
}

// UpdateSwipe updates a swipe in the database.
func (s swipeRepository) UpdateSwipe(ctx context.Context, payload *models.Swipe) (*models.Swipe, error) {
	result, err := s.mongo.coll(s.collection).UpdateOne(
		ctx,
		bson.M{"id": payload.ID},
		bson.M{"$set": payload},
	)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no swipe updated")
	}
	return payload, nil
}

// DeleteSwipe deletes a swipe from the database.
func (s swipeRepository) DeleteSwipe(ctx context.Context, id string) error {
	result, err := s.mongo.coll(s.collection).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no swipe deleted")
	}
	return nil
}

func NewSwipeRepo(store *MongoStore) repository.SwipesRepository {
	return &swipeRepository{
		mongo:          store,
		collection:     constants.SwipeCollection,
		userCollection: constants.UserCollection,
	}
}
