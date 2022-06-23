package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"testing"
	"time"
)

var (
	ctx         context.Context
	constructor *handler
	srv         *MockService
)

type MockService struct {
	mock.Mock
}

func (m *MockService) ListPosts(in *sma.ListPostsReq) (*sma.ListPostsResp, error) {
	args := m.Called(in)
	return args.Get(0).(*sma.ListPostsResp), args.Error(1)
}

func (m *MockService) RemovePost(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockService) RetrievePost(id primitive.ObjectID) (*sma.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*sma.Post), args.Error(1)
}

func (m *MockService) UpdatePost(id primitive.ObjectID, in *sma.UpdatePost) error {
	args := m.Called(id, in)
	return args.Error(0)
}

func (m *MockService) CreatePost(in *sma.CreatePostReq) (*sma.Post, error) {
	args := m.Called(in)
	return args.Get(0).(*sma.Post), args.Error(1)
}

type protoDef struct {
	empty bool
	value *sma.Post
}

type expectDef struct {
	value    interface{}
	typeName string
}

func (def *protoDef) GetValue() *sma.Post {
	return def.value
}

func (def *protoDef) SetValue() {
	if def.value != nil {
		return
	}

	if def.empty {
		def.value = &sma.Post{}
		return
	}

	now := time.Now()
	def.value = &sma.Post{
		Id:          primitive.NewObjectID().Hex(),
		Title:       "ExampleMock",
		Description: "Mock description one",
		UserId:      uuid.New().String(),
		Status:      sma.PostStatus_POST_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
}

func (def *protoDef) Clear() {
	if def.value != nil {
		def.value = nil
	}

	def.empty = false
}

func NewProtoDef(empty bool) *protoDef {
	return &protoDef{empty: empty}
}

func (e *expectDef) Clear() {
	if e.value != nil {
		e.value = nil
	}

	e.typeName = ""
}

func (e *expectDef) IsError() bool {
	return strings.EqualFold(e.typeName, "error")
}

func (e *expectDef) SetValue(payload interface{}) {
	if val, ok := e.value.(error); ok {
		e.value = val
		e.typeName = "error"
		return
	}

	e.value = payload
	e.typeName = "model"
}

func (e *expectDef) GetValue() interface{} {
	return e.value
}

func NewExpected(v interface{}) *expectDef {
	return &expectDef{value: v}
}

func init() {
	if constructor != nil && ctx != nil && srv != nil {
		return
	}

	srv = new(MockService)
	constructor = New(srv)
	ctx = context.Background()
}

func eval(t *testing.T, ok bool) {
	if !ok {
		t.FailNow()
	}
}
