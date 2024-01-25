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

type matchRepository struct {
	mongo      *MongoStore
	collection string
}

// GetMatchByProfiles returns a match by the given profiles.
func (m matchRepository) GetMatchByProfiles(ctx context.Context, profiles []string) (*models.Match, error) {
	var match models.Match
	err := m.mongo.coll(m.collection).FindOne(ctx, bson.M{"profiles": profiles}).Decode(&match)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

// CreateMatch creates a new match in the database.
func (m matchRepository) CreateMatch(ctx context.Context, payload *models.Match) (*models.Match, error) {
	existingMatch, err := m.GetMatchByProfiles(ctx, payload.Profiles)
	if err != nil {
		return nil, err
	}

	if existingMatch != nil {
		return nil, repository.ErrDuplicateFound
	}

	_, err = m.mongo.coll(m.collection).InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// GetMatchById returns a match by their ID.
func (m matchRepository) GetMatchById(ctx context.Context, id string) (*models.Match, error) {
	var match models.Match
	err := m.mongo.coll(m.collection).FindOne(ctx, bson.M{"_id": id}).Decode(&match)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

// GetMatchesFiltered returns a match by the given filter.
func (m matchRepository) GetMatchesFiltered(ctx context.Context, filter models.MatchFilter) (*models.MatchedUserInfo, error) {
	qb := NewMatchedUserInfoQueryBuilder(filter).
		MatchProfiles([]string{filter.UserID}).
		UnwindProfiles().
		LookupUsers().
		UnwindUsers().
		GroupResults()

	cursor, err := m.mongo.coll(m.collection).Aggregate(ctx, qb.Build())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := &models.MatchedUserInfo{}
	if cursor.Next(ctx) {
		if err := cursor.Decode(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// UpdateMatch updates a match in the database.
func (m matchRepository) UpdateMatch(ctx context.Context, payload *models.Match) (*models.Match, error) {
	update := bson.M{"$set": payload}
	result, err := m.mongo.coll(m.collection).UpdateOne(
		ctx,
		bson.M{"id": payload.ID},
		update,
	)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("no matching document found")
	}
	return payload, nil
}

// DeleteMatch deletes a match in the database.
func (m matchRepository) DeleteMatch(ctx context.Context, id string) error {
	result, err := m.mongo.coll(m.collection).DeleteOne(
		ctx,
		bson.M{"id": id},
	)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no matching document found")
	}
	return nil
}

func NewMatchRepo(store *MongoStore) repository.MatchRepository {
	return &matchRepository{
		mongo:      store,
		collection: constants.MatchCollection,
	}
}
