package mongodb

import (
	"api/constants"
	"api/models"
	"api/repository"
	"api/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type userRepository struct {
	mongo      *MongoStore
	collection string
	jwtSecret  string
}

// Discover returns a list of users that match the given filter.
func (u userRepository) Discover(ctx context.Context, filter models.UserFilter, user models.User) ([]*models.User, error) {
	qb := NewDiscoverQueryBuilder(user, filter).
		LookupSwipes().
		MatchSwipesEmpty().
		LookupSwipesCount().
		Projection().
		AgeFilter(filter.MinAge, filter.MaxAge)

	cursor, err := u.mongo.coll(u.collection).Aggregate(ctx, qb.Build())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var result []*models.User
	if err := cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserCount returns the total number of users in the database.
func (u userRepository) GetUserCount(ctx context.Context) (int, error) {
	count, err := u.mongo.coll(u.collection).CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// InsertUsers inserts a list of users into the database.
func (u userRepository) InsertUsers(ctx context.Context, users []*models.User) error {
	var documents []interface{}
	for _, user := range users {
		user.ID = utils.GenerateId()
		documents = append(documents, user)
	}
	_, err := u.mongo.coll(u.collection).InsertMany(ctx, documents)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser creates a new user in the database.
func (u userRepository) CreateUser(ctx context.Context, payload *models.User) (*models.User, error) {

	filters := bson.M{
		"$or": []bson.M{
			{
				"email": payload.Email,
			},
			// if we want to add more fields to check for duplicity, we can add them here.
		},
	}

	var Profile models.User
	if err := u.mongo.coll(u.collection).FindOne(ctx, filters).Decode(&Profile); err == nil {
		return nil, repository.ErrDuplicateFound
	}
	payload.ID = utils.GenerateId()
	_, err := u.mongo.coll(u.collection).InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// GetUserById returns a user by their ID.
func (u userRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return u.GetProfileByField(ctx, "id", id)
}

// GetUserByEmail returns a user by their email.
func (u userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.GetProfileByField(ctx, "email", email)
}

// UpdateUser updates a user in the database.
func (u userRepository) UpdateUser(ctx context.Context, payload *models.User) (*models.User, error) {
	update := bson.M{"$set": payload}
	_, err := u.mongo.coll(u.collection).UpdateOne(ctx, bson.M{"id": payload.ID}, update)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// GetProfileByField returns a user by the given field.
func (u userRepository) GetProfileByField(ctx context.Context, field, value string) (*models.User, error) {
	var Profile models.User
	if err := u.mongo.coll(u.collection).FindOne(ctx, bson.M{field: value}).Decode(&Profile); err != nil {
		return nil, err
	}
	return &Profile, nil
}

func NewUserRepo(store *MongoStore) repository.UserRepository {
	return &userRepository{
		mongo:      store,
		collection: constants.UserCollection,
	}
}
