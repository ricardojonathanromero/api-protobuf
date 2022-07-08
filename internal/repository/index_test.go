package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
	"time"
)

type MockDB struct{ mock.Mock }

func (m *MockDB) Connect() (*mongo.Client, error) {
	args := m.Called()
	return args.Get(0).(*mongo.Client), args.Error(1)
}

func (m *MockDB) Close() {}

func TestRepository_CountDocuments(t *testing.T) {
	tests := []struct {
		title           string
		postFind        *sma.ListPostsReq
		errorCollection error
		isSuccess       bool
		isError         bool
		expected        interface{}
	}{
		{
			title: "happy_path",
			postFind: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
			},
			isSuccess: true,
		},
		{
			title: "happy_path_filter",
			postFind: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
				Filter:  sma.Filters_FILTER_ACTIVE,
			},
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			postFind:        &sma.ListPostsReq{},
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				// return two mock responses, one for cursor and second for find response
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "posts.posts", mtest.FirstBatch, bson.D{
					{Key: "n", Value: 1},
				}))
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    2,
					Message: "no exist documents",
				}))
			}

			count, err := repo.CountDocuments(test.postFind)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
				mt.Log(count)
			}
		})
	}
}

func TestRepository_DeleteDocumentByID(t *testing.T) {
	tests := []struct {
		title           string
		id              primitive.ObjectID
		errorCollection error
		isSuccess       bool
		isError         bool
		isErrorDelete   bool
		expected        interface{}
	}{
		{
			title:     "happy_path",
			id:        primitive.NewObjectID(),
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			id:              primitive.NewObjectID(),
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
		{
			title:    "error_deleting_document",
			id:       primitive.NewObjectID(),
			isError:  true,
			expected: errors.New("mongo: document not deleted"),
		},
		{
			title:         "error_document_not_deleted",
			id:            primitive.NewObjectID(),
			isError:       true,
			isErrorDelete: true,
			expected:      mongo.ErrNoDocuments,
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "n", Value: 1},
				})
			} else if test.isErrorDelete {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "n", Value: 0},
				})
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    4,
					Message: "document not deleted",
				}))
			}

			err := repo.DeleteDocumentByID(test.id)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
			}
		})
	}
}

func TestRepository_FindDocuments(t *testing.T) {
	tests := []struct {
		title           string
		postFind        *sma.ListPostsReq
		errorCollection error
		isSuccess       bool
		isError         bool
		expected        interface{}
	}{
		{
			title: "happy_path",
			postFind: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
			},
			isSuccess: true,
		},
		{
			title: "happy_path_filter",
			postFind: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
				Filter:  sma.Filters_FILTER_ACTIVE,
			},
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			postFind:        &sma.ListPostsReq{},
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
		{
			title: "error_find_result",
			postFind: &sma.ListPostsReq{
				UserId:  uuid.New().String(),
				Page:    1,
				PerPage: 10,
				Filter:  sma.Filters_FILTER_DRAFT,
			},
			isError:  true,
			expected: errors.New("no documents found"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				// return two mock responses, one for cursor and second for find response
				first := mtest.CreateCursorResponse(1, "posts.posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "title", Value: "example"},
					{Key: "description", Value: "example"},
					{Key: "user_id", Value: test.postFind.UserId},
					{Key: "status", Value: sma.PostStatus_POST_STATUS_ACTIVE.String()},
					{Key: "created_at", Value: time.Now()},
					{Key: "update_at", Value: time.Now()},
				})
				second := mtest.CreateCursorResponse(1, "posts.posts", mtest.NextBatch, bson.D{
					{Key: "_id", Value: primitive.NewObjectID()},
					{Key: "title", Value: "example_two"},
					{Key: "description", Value: "example two"},
					{Key: "user_id", Value: test.postFind.UserId},
					{Key: "status", Value: sma.PostStatus_POST_STATUS_ACTIVE.String()},
					{Key: "created_at", Value: time.Now()},
					{Key: "update_at", Value: time.Now()},
				})
				killCursors := mtest.CreateCursorResponse(0, "posts.posts", mtest.NextBatch)

				mt.AddMockResponses(first, second, killCursors)
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    2,
					Message: "document not found",
				}))
			}

			postList, err := repo.FindDocuments(test.postFind)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
				mt.Log(postList)
			}
		})
	}
}

func TestRepository_FindDocumentByID(t *testing.T) {
	tests := []struct {
		title           string
		id              primitive.ObjectID
		errorCollection error
		isSuccess       bool
		isError         bool
		expected        interface{}
	}{
		{
			title:     "happy_path",
			id:        primitive.NewObjectID(),
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			id:              primitive.NewObjectID(),
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
		{
			title:    "error_find_result",
			id:       primitive.NewObjectID(),
			isError:  true,
			expected: errors.New("document not found"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				// return two mock responses, one for cursor and second for find response
				first := mtest.CreateCursorResponse(1, "posts.posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: test.id},
					{Key: "title", Value: "example"},
					{Key: "description", Value: "example"},
					{Key: "user_id", Value: uuid.New().String()},
					{Key: "status", Value: sma.PostStatus_POST_STATUS_ACTIVE.String()},
					{Key: "created_at", Value: time.Now()},
					{Key: "update_at", Value: time.Now()},
				})
				mt.AddMockResponses(first)
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    2,
					Message: "document not found",
				}))
			}

			post, err := repo.FindDocumentByID(test.id)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
				eval(mt, assert.Equal(mt, post.ID, test.id))
				mt.Log(post)
			}
		})
	}
}

func TestRepository_InsertDocument(t *testing.T) {
	now := time.Now()

	tests := []struct {
		title           string
		doc             *models.Post
		errorCollection error
		isSuccess       bool
		isError         bool
		expected        interface{}
	}{
		{
			title: "happy_path",
			doc: &models.Post{
				ID:          primitive.ObjectID{},
				Title:       "Happy path example",
				Description: "example",
				UserID:      uuid.New().String(),
				MediaIDs:    []string{"hello"},
				Status:      "active",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			doc:             &models.Post{},
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
		{
			title:    "error_inserting_document",
			doc:      &models.Post{},
			isError:  true,
			expected: errors.New("document not saved"),
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    11000,
					Message: "document not saved",
				}))
			}

			id, err := repo.InsertDocument(test.doc)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
				eval(mt, primitive.IsValidObjectID(id.Hex()))
			}
		})
	}
}

func TestRepository_UpdateDocumentByID(t *testing.T) {
	tests := []struct {
		title           string
		id              primitive.ObjectID
		doc             *models.Post
		errorCollection error
		isSuccess       bool
		isError         bool
		isErrorUpdate   bool
		expected        interface{}
	}{
		{
			title: "happy_path",
			id:    primitive.NewObjectID(),
			doc: &models.Post{
				Title: "example",
			},
			isSuccess: true,
		},
		{
			title:           "error_mongo_client",
			id:              primitive.NewObjectID(),
			doc:             &models.Post{},
			errorCollection: errors.New("error connecting to mongo client"),
			isError:         true,
			expected:        errors.New("error connecting to mongo client"),
		},
		{
			title:    "error_updating_document",
			id:       primitive.NewObjectID(),
			doc:      &models.Post{},
			isError:  true,
			expected: errors.New("mongo: document not updated"),
		},
		{
			title:         "error_document_not_updated",
			id:            primitive.NewObjectID(),
			doc:           &models.Post{},
			isError:       true,
			isErrorUpdate: true,
			expected:      mongo.ErrNoDocuments,
		},
	}

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, test := range tests {
		mt.Run(test.title, func(mt *mtest.T) {
			mockDB := new(MockDB)
			// evaluates db.IDB instance mock to return error or *mongo.Client
			mockDBConn(mockDB, mt.Client, test.errorCollection)

			repo := New(mockDB)

			if test.isSuccess {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 1},
					{Key: "value", Value: bson.D{
						{Key: "_id", Value: test.id},
					}},
				})
			} else if test.isErrorUpdate {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "acknowledged", Value: true},
					{Key: "n", Value: 0},
				})
			} else {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   1,
					Code:    4,
					Message: "document not deleted",
				}))
			}

			err := repo.UpdateDocumentByID(test.id, test.doc)
			mt.Log(err)
			if test.isError {
				eval(mt, assert.Error(mt, err))
				eval(mt, assert.Equal(mt, err, test.expected))
			} else {
				eval(mt, assert.NoError(mt, err))
			}
		})
	}
}

func mockDBConn(m *MockDB, client *mongo.Client, err error) {
	if err != nil {
		m.On("Connect").Return(client, err)
	} else {
		m.On("Connect").Return(client, nil)
	}
}

func eval(t *mtest.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}
