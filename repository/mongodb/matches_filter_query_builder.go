package mongodb

import (
	"api/constants"
	"api/models"
	"go.mongodb.org/mongo-driver/bson"
)

// MatchedUserInfoQueryBuilder is a builder for constructing MongoDB aggregation pipelines
// specifically for matching user information.
type MatchedUserInfoQueryBuilder struct {
	filter   models.MatchFilter
	pipeline []bson.M
}

// NewMatchedUserInfoQueryBuilder creates and returns a new instance of MatchedUserInfoQueryBuilder.
func NewMatchedUserInfoQueryBuilder(filter models.MatchFilter) *MatchedUserInfoQueryBuilder {
	return &MatchedUserInfoQueryBuilder{
		filter:   filter,
		pipeline: []bson.M{},
	}
}

func (qb *MatchedUserInfoQueryBuilder) MatchProfiles(profiles []string) *MatchedUserInfoQueryBuilder {
	matchFilter := bson.M{"profiles": bson.M{"$in": profiles}}
	qb.pipeline = append(qb.pipeline, bson.M{"$match": matchFilter})
	return qb
}

func (qb *MatchedUserInfoQueryBuilder) UnwindProfiles() *MatchedUserInfoQueryBuilder {
	qb.pipeline = append(qb.pipeline, bson.M{"$unwind": "$profiles"})
	return qb
}

func (qb *MatchedUserInfoQueryBuilder) LookupUsers() *MatchedUserInfoQueryBuilder {
	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":         constants.UserCollection,
			"localField":   "profiles",
			"foreignField": "id",
			"as":           "user",
		},
	}
	qb.pipeline = append(qb.pipeline, lookupStage)
	return qb
}

func (qb *MatchedUserInfoQueryBuilder) UnwindUsers() *MatchedUserInfoQueryBuilder {
	qb.pipeline = append(qb.pipeline, bson.M{"$unwind": "$user"})
	return qb
}

func (qb *MatchedUserInfoQueryBuilder) GroupResults() *MatchedUserInfoQueryBuilder {
	groupStage := bson.M{
		"$group": bson.M{
			"_id":           nil,
			"current_user":  bson.M{"$first": "$user"},
			"matched_users": bson.M{"$push": "$user"},
		},
	}
	qb.pipeline = append(qb.pipeline, groupStage)
	return qb
}

func (qb *MatchedUserInfoQueryBuilder) Build() []bson.M {
	return qb.pipeline
}
