package db

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// All tests that use mtest.Setup() are expected to be integration tests, so skip them when the
	// -short flag is included in the "go test" command. Also, we have to parse flags here to use
	// testing.Short() because flags aren't parsed before TestMain() is called.
	flag.Parse()

	log.Info("init test main")
	if testing.Short() {
		log.Info("skipping mongo-test integration test in short mode")
		return
	}

	log.Info("setup mtest")
	if err := mtest.Setup(); err != nil {
		log.Fatalf("error setup mtest: %v", err)
	}

	defer os.Exit(m.Run())
	if err := mtest.Teardown(); err != nil {
		log.Fatalf("error teardown mtest: %v", err)
	}
}

func TestDBConnection(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.AddMockResponses(bson.D{{"ok", 1}})

	log.Info("init constructor")
	constructor := New()

	_ = os.Setenv("MONGO_URI", mtest.ClusterURI())
	conn, err := constructor.Connect()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
		t.FailNow()
	}

	t.Log(conn)
}
