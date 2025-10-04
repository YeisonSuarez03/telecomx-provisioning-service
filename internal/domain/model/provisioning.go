package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Provisioning struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"userId" json:"userId"`
	ServiceName   string             `bson:"serviceName" json:"serviceName"`
	Status        string             `bson:"status" json:"status"`
	LastUpdatedAt string             `bson:"lastUpdatedAt" json:"lastUpdatedAt"`
	CreatedAt     string             `bson:"createdAt" json:"createdAt"`
}
