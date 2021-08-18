package mongo

import (
	"context"
	"errors"

	"github.com/buraksekili/getir-task/pkg"

	"github.com/buraksekili/getir-task/persistence"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrDecodeData represents a failure during data decode.
	ErrDecodeData = errors.New("failed to decode a data from mongodb")

	// ErrParseTime represents a failure during time parsing.
	ErrParseTime = errors.New("failed to parse time")
)

// repository represents a new mongo repository.
type repository struct {
	collection *mongo.Collection
	client     *mongo.Client
}

// NewMongo returns a new MongoDB.
func NewMongo(collection *mongo.Collection) persistence.Database {
	return repository{collection: collection}

}

func (m repository) FetchData(ctx context.Context, f persistence.Filter) ([]persistence.FetchResObj, error) {
	t, err := pkg.ParseTime(f.StartDate)
	if err != nil {
		return []persistence.FetchResObj{}, pkg.Wrap(ErrParseTime, err)
	}

	t1, err := pkg.ParseTime(f.EndDate)
	if err != nil {
		return []persistence.FetchResObj{}, pkg.Wrap(ErrParseTime, err)
	}

	pipe := []bson.M{{"$project": bson.M{
		"key":        1,
		"createdAt":  1,
		"totalCount": bson.M{"$sum": "$counts"},
	}}, {"$match": bson.M{
		"$and": []bson.M{
			{"createdAt": bson.M{"$gt": t}},
			{"createdAt": bson.M{"$lt": t1}},
			{"totalCount": bson.M{"$gt": f.MinTotalCount}},
			{"totalCount": bson.M{"$lt": f.MaxTotalCount}}},
	}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipe)
	if err != nil {
		return []persistence.FetchResObj{}, pkg.Wrap(persistence.ErrAggregation, err)
	}
	defer cursor.Close(ctx)

	r := []persistence.FetchResObj{}
	for cursor.Next(context.Background()) {
		o := persistence.FetchResObj{}
		if err := cursor.Decode(&o); err != nil {
			return []persistence.FetchResObj{}, pkg.Wrap(ErrDecodeData, err)
		}
		r = append(r, o)
	}
	return r, nil
}
