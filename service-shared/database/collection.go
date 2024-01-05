//Package database provides functions for interacting with a persistent db
package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoCaller represents operations called on a mongo collection. It is a wrapper interface to aid testing.
type MongoCaller interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type MongoCollection struct {
	collection *mongo.Collection
}

func NewMongoCollection(collection *mongo.Collection) MongoCollection {
	return MongoCollection{collection: collection}
}

func (m MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(ctx, document, opts...)
}

func (m MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.collection.FindOne(ctx, filter, opts...)
}

func (m MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	return m.collection.Find(ctx, filter, opts...)
}

func (m MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.collection.DeleteOne(ctx, filter, opts...)
}

func (m MongoCollection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.collection.UpdateByID(ctx, id, update, opts...)
}
