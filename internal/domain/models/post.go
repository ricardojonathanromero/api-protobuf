package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	UserID      string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	MediaIDs    []string           `json:"media_ids,omitempty" bson:"media_ids,omitempty"`
	Status      string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type PostList struct {
	Posts []*Post `json:"posts" bson:"posts,omitempty"`
	Count int64   `json:"count" bson:"count,omitempty"`
}
