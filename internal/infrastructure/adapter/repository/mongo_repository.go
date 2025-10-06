package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"telecomx-provisioning-service/internal/domain/model"
)

type MongoRepository struct {
	Collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{Collection: db.Collection("provisioning")}
}

func (r *MongoRepository) Create(ctx context.Context, p *model.Provisioning) error {
	p.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.Collection.InsertOne(ctx, p)
	return err
}

func (r *MongoRepository) UpdateStatus(ctx context.Context, userID, status string) error {
	filter := bson.M{"userId": userID}
	update := bson.M{"$set": bson.M{"status": status, "lastUpdatedAt": time.Now().Format(time.RFC3339)}}
	_, err := r.Collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *MongoRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.Collection.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]model.Provisioning, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Provisioning
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
