package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func Test(t *testing.T) {
	now := timestamppb.Now().AsTime()
	post := Post{
		ID:          primitive.NewObjectID(),
		Title:       "Example",
		Description: "Example Description",
		UserID:      uuid.New().String(),
		MediaIDs:    []string{"one"},
		Status:      "active",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	postM, _ := json.Marshal(post)
	t.Log(string(postM))

	list := PostList{
		Posts: []*Post{&post},
		Count: 1,
	}

	postListM, _ := json.Marshal(list)
	t.Log(string(postListM))
}
