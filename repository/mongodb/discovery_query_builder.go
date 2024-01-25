package mongodb

import (
	"api/constants"
	"api/models"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// DiscoverQueryBuilder is a builder for constructing MongoDB aggregation pipelines
// for the "Discover" feature, which likely involves finding and filtering users.
type DiscoverQueryBuilder struct {
	stages []bson.M
}

// NewDiscoverQueryBuilder initializes a new DiscoverQueryBuilder with a geospatial search stage.
func NewDiscoverQueryBuilder(user models.User, filter models.UserFilter) *DiscoverQueryBuilder {
	geoNearStage := bson.M{
		"$geoNear": bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": user.Location,
			},
			"distanceField": "distance",
			"spherical":     true,
		},
	}

	// If a maximum distance filter is set, apply it.
	if filter.MaxDistance > 0 {
		geoNearStage["$geoNear"].(bson.M)["maxDistance"] = filter.MaxDistance * 1000
	}
	return &DiscoverQueryBuilder{stages: []bson.M{geoNearStage}}
}

// LookupSwipes adds a stage to the pipeline to look up swipe data for each user.
func (qb *DiscoverQueryBuilder) LookupSwipes() *DiscoverQueryBuilder {
	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":     constants.SwipeCollection,
			"let":      bson.M{"targetUserId": "$id"},
			"pipeline": []bson.M{{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$prospect_id", "$$targetUserId"}}}}},
			"as":       "prospects",
		},
	}
	qb.stages = append(qb.stages, lookupStage)
	return qb
}

// MatchSwipesEmpty filters out users who have already been swiped on.
func (qb *DiscoverQueryBuilder) MatchSwipesEmpty() *DiscoverQueryBuilder {
	matchStage := bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				{"id": bson.M{"$ne": "user.ID"}},
				{"prospects": bson.M{"$eq": []interface{}{}}},
			},
		},
	}
	qb.stages = append(qb.stages, matchStage)
	return qb
}

// LookupSwipesCount adds a stage to count the number of swipes for each user.
func (qb *DiscoverQueryBuilder) LookupSwipesCount() *DiscoverQueryBuilder {
	swipesLookupStage := bson.M{
		"$lookup": bson.M{
			"from": "swipes",
			"let":  bson.M{"targetUserId": "$id"},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$user_id", "$$targetUserId"}},
								{"$eq": []interface{}{"$prospect_id", "$$targetUserId"}},
							},
						},
					},
				},
			},
			"as": constants.SwipeCollection,
		},
	}
	qb.stages = append(qb.stages, swipesLookupStage, bson.M{"$match": bson.M{"swipes": bson.M{"$eq": []interface{}{}}}}, bson.M{"$match": bson.M{"id": bson.M{"$ne": "user.ID"}}})
	return qb
}

func (qb *DiscoverQueryBuilder) Projection() *DiscoverQueryBuilder {
	projectionStage := bson.M{
		"$project": bson.M{
			"id":                1,
			"name":              1,
			"gender":            1,
			"age":               bson.M{"$toInt": bson.M{"$divide": []interface{}{bson.M{"$subtract": []interface{}{time.Now(), bson.M{"$toDate": "$date_of_birth"}}}, 31556952000}}},
			"date_of_birth":     1,
			"location":          1,
			"height":            1,
			"ethnicity":         1,
			"pets":              1,
			"sexuality":         1,
			"religion":          1,
			"drinking":          1,
			"smoking":           1,
			"drugs":             1,
			"dating_intentions": 1,
			"kids":              1,
			"occupation":        1,
			"swipe_count":       1,
			"attractiveness":    1,
			"bio":               1,
			"distance":          bson.M{"$divide": []interface{}{"$distance", 1000}},
		},
	}
	qb.stages = append(qb.stages, projectionStage)
	return qb
}

// AgeFilter adds a filtering stage based on age range.
func (qb *DiscoverQueryBuilder) AgeFilter(minAge, maxAge int) *DiscoverQueryBuilder {
	if minAge > 0 || maxAge > 0 {
		ageFilterStage := bson.M{"$match": bson.M{"age": bson.M{"$gte": minAge, "$lte": maxAge}}}
		qb.stages = append(qb.stages, ageFilterStage)
	}
	return qb
}

func (qb *DiscoverQueryBuilder) Build() []bson.M {
	return qb.stages
}
