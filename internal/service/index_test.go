package service

import (
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type MockRepo struct{ mock.Mock }
type MockUtils struct{ mock.Mock }

func (m *MockRepo) CountDocuments(in *sma.ListPostsReq) (int64, error) {
	args := m.Called(in)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) DeleteDocumentByID(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) FindDocumentByID(id primitive.ObjectID) (*models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *MockRepo) FindDocuments(input *sma.ListPostsReq) ([]*models.Post, error) {
	args := m.Called(input)
	return args.Get(0).([]*models.Post), args.Error(1)
}

func (m *MockRepo) InsertDocument(doc interface{}) (primitive.ObjectID, error) {
	args := m.Called(doc)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockRepo) UpdateDocumentByID(id primitive.ObjectID, doc interface{}) error {
	args := m.Called(id, doc)
	return args.Error(0)
}

func (m *MockUtils) Now() *time.Time {
	args := m.Called()
	return args.Get(0).(*time.Time)
}

func eval(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}
