package protobuf

import (
	"github.com/ricardojonathanromero/api-protobuf/infrastructure/db"
	"github.com/ricardojonathanromero/api-protobuf/internal/handler"
	"github.com/ricardojonathanromero/api-protobuf/internal/repository"
	"github.com/ricardojonathanromero/api-protobuf/internal/service"
	"github.com/ricardojonathanromero/api-protobuf/proto/sma"
	"github.com/ricardojonathanromero/api-protobuf/utils"
	"google.golang.org/grpc"
	"net"
)

type IProtobuf interface {
	Server(gRPCAddr string, dbInstance db.IDB) error
}

type protobuf struct{}

var _ IProtobuf = (*protobuf)(nil)

func (pb *protobuf) Server(gRPCAddr string, dbInstance db.IDB) error {
	// create new gRPC server
	srv := grpc.NewServer()

	// creating handlers
	repo := repository.New(dbInstance)
	u := utils.New()
	serv := service.New(repo, u)

	// register the GreeterServerImpl on the gRPC server
	sma.RegisterPostsServer(srv, handler.New(serv))

	// start listening on port :8080 for a tcp connection
	l, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		return err
	}

	// the start gRPC server
	return srv.Serve(l)
}

func New() IProtobuf {
	return &protobuf{}
}
