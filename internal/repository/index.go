package repository

import (
	"errors"
	"github.com/ricardojonathanromero/api-protobuf/infrastructure/db"
	"github.com/ricardojonathanromero/api-protobuf/internal/domain/models"
	"github.com/ricardojonathanromero/api-protobuf/internal/port"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/ricardojonathanromero/api-protobuf/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	instance db.IDB
}

var _ port.IRepository = (*repository)(nil)

const (
	deadline   int    = 10
	dbName     string = "posts"
	collection string = "posts"
)

func (repo *repository) CountDocuments(input *sma.ListPostsReq) (int64, error) {
	col, err := repo.collection()
	if err != nil {
		return 0, err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	filter := bson.D{{"user_id", input.UserId}}

	if input.Filter > 0 {
		filter = append(filter, bson.E{Key: "status", Value: input.Filter})
	}

	return col.CountDocuments(ctx, filter)
}

// DeleteDocumentByID removes document from db using id/*
func (repo *repository) DeleteDocumentByID(id primitive.ObjectID) error {
	// removing document from db
	col, err := repo.collection()
	if err != nil {
		return err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	res, err := col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("mongo: document not deleted")
	}

	if res.DeletedCount <= 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// FindDocumentByID find documents by id/*
func (repo *repository) FindDocumentByID(id primitive.ObjectID) (*models.Post, error) {
	var result *models.Post

	col, err := repo.collection()
	if err != nil {
		return result, err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	err = col.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return result, errors.New("document not found")
	}

	return result, nil
}

// FindDocuments finds documents/*
func (repo *repository) FindDocuments(input *sma.ListPostsReq) ([]*models.Post, error) {
	result := make([]*models.Post, 0)

	col, err := repo.collection()
	if err != nil {
		return result, err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	filter := bson.D{{"user_id", input.UserId}}
	opts := options.
		Find().
		SetLimit(input.PerPage).
		SetSkip((input.Page - 1) * input.PerPage)

	if input.Filter > 0 {
		filter = append(filter, bson.E{Key: "status", Value: input.Filter})
	}

	// find documents
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return result, errors.New("no documents found")
	}

	_ = cursor.All(ctx, &result)

	return result, nil
}

// InsertDocument creates a new document into collection.
// requires: interface{}
// produces: primitive.ObjectID or error/*
func (repo *repository) InsertDocument(doc interface{}) (primitive.ObjectID, error) {
	col, err := repo.collection()
	if err != nil {
		return primitive.NilObjectID, err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	// inserting document received
	res, err := col.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, errors.New("document not saved")
	}

	// returning primitive id doing cast to result
	return res.InsertedID.(primitive.ObjectID), nil
}

// UpdateDocumentByID updates document information filtering by _id/*
func (repo *repository) UpdateDocumentByID(id primitive.ObjectID, doc interface{}) error {
	col, err := repo.collection()
	if err != nil {
		return err
	}

	ctx, cancel := utils.ContextWithTimeout(deadline)
	defer cancel()

	res, err := col.UpdateByID(ctx, id, bson.M{"$set": doc})
	if err != nil {
		return errors.New("mongo: document not updated")
	}

	if res.ModifiedCount <= 0 && res.MatchedCount <= 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// collection returns *mongo.Collection pointer/*
func (repo *repository) collection() (*mongo.Collection, error) {
	var result *mongo.Collection

	client, err := repo.instance.Connect()
	if err != nil {
		return result, err
	}

	return client.Database(dbName).Collection(collection), nil
}

// New constructor function/*
func New(instance db.IDB) port.IRepository {
	return &repository{instance: instance}
}
