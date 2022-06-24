package db

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	cl "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func createMongoServer() (*cl.Client, string) {
	cli, err := cl.NewClientWithOpts()
	if err != nil {
		log.Fatalf("error creating docker client. %v", err)
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "mongo",
		ExposedPorts: nat.PortSet{"27017": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{"27017": {{HostIP: "localhost", HostPort: "27017"}}},
	}, &network.NetworkingConfig{}, &v1.Platform{}, "mongo-test")
	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("error starting container: %v", err)
	}

	time.Sleep(2 * time.Second)
	log.Info("container configured")
	return cli, resp.ID
}

func TestDBConnection_ErrorPingConnection(t *testing.T) {
	log.Info("init constructor")
	constructor := New()

	_ = os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer func() {
		_ = os.Remove("MONGO_URI")
	}()

	_, err := constructor.Connect()
	if !assert.Error(t, err) {
		t.Error("error expected at this point")
		t.FailNow()
	}

	t.Log("success")
}

func TestDBConnection_ErrorConnection(t *testing.T) {
	log.Info("init constructor")
	constructor := New()

	_, err := constructor.Connect()
	if !assert.Error(t, err) {
		t.Error("error expected at this point")
		t.FailNow()
	}

	t.Log("success")
}

func TestDBConnection_HappyPath(t *testing.T) {
	cli, id := createMongoServer()

	defer func() {
		if cli != nil {
			ctx := context.Background()
			_ = cli.ContainerStop(ctx, id, nil)
			_ = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true})
		}
		_ = os.Remove("MONGO_URI")
	}()

	log.Info("init constructor")
	constructor := New()
	defer constructor.Close()

	_ = os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	_, err := constructor.Connect()
	if !assert.NoError(t, err) {
		t.Errorf("error connecting to db: %v", err)
		t.FailNow()
	}

	_, err = constructor.Connect()
	if !assert.NoError(t, err) {
		t.Errorf("error connecting to db: %v", err)
		t.FailNow()
	}

	t.Log("success")
}
