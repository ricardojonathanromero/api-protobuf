package db

import (
	"context"
	"fmt"
	"github.com/ricardojonathanromero/api-protobuf/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

type IDB interface {
	Close()
	Connect() (*mongo.Client, error)
}

type conn struct{}

var (
	_      IDB = (*conn)(nil)
	once   sync.Once
	client *mongo.Client
)

// Close disconnect db client connection in case an
// active connection exists./*
func (c *conn) Close() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Errorf("error closing mongo client. reason => %v", err)
		}
	}
}

// Connect creates or reuse mongo connection, or returns
// an error instead./*
func (c *conn) Connect() (*mongo.Client, error) {
	if client != nil {
		fmt.Println("skipping constructor")
		return client, nil
	}

	var err error
	once.Do(func() {
		// get environment variable
		uri := utils.GetEnv("MONGO_URI", "")

		// create client connection
		log.Infof("conencting to %s", uri)
		client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			log.Errorf("error connecting to db. reason\n%v", err)
			return
		}

		// create context with timeout to validate db connection
		ctx, cancel := utils.ContextWithTimeout(10)
		defer cancel()

		// confirm connection making ping to db server
		log.Info("doing ping to db")
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			_ = client.Disconnect(context.Background())
			log.Errorf("error ping to db. reason \n%v", err)
			return
		}
	})

	return client, err
}

// New constructor function /*
func New() IDB {
	return &conn{}
}
