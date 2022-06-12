package port

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	DeleteDocumentByID(id primitive.ObjectID) error
	FindDocumentByID(id primitive.ObjectID) (*models.Post, error)
	FindDocuments(input *sma.ListPostsReq) (*models.PostList, error)
	InsertDocument(doc interface{}) (primitive.ObjectID, error)
	UpdateDocumentByID(id primitive.ObjectID, doc interface{}) error
}

type IService interface {
	CreatePost(in *sma.CreatePostReq) (*sma.Post, error)
	ListPosts(in *sma.ListPostsReq) (*sma.ListPostsResp, error)
	RemovePost(id primitive.ObjectID) error
	RetrievePost(id primitive.ObjectID) (*sma.Post, error)
	UpdatePost(id primitive.ObjectID, in *sma.UpdatePost) error
}
