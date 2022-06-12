package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/ricardojonathanromero/api-protobuf/infrastructure/db"
	"github.com/ricardojonathanromero/api-protobuf/infrastructure/protobuf"
	"github.com/ricardojonathanromero/api-protobuf/infrastructure/routes"
	log "github.com/sirupsen/logrus"
)

var (
	enableProxy = flag.Bool("enable_proxy", false, "flag to configure http proxy server")
	grpcAddr    = flag.String("grpc_addr", ":8090", "endpoint of the gRPC server")
	httpAddr    = flag.String("http_addr", ":8080", "endpoint of the http server")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// creating new db connection
	dbInstance := db.New()

	// defer function used to close client in case a panic error occurs
	defer dbInstance.Close()

	if *enableProxy {
		log.Info("creating http proxy server")
		go func() {
			if err := routes.New().ProxyServer(*grpcAddr, *httpAddr); err != nil {
				log.Fatal(err)
			}
		}()
	}

	log.Info("creating gRPC server")
	if err := protobuf.New().Server(*grpcAddr, dbInstance); err != nil {
		log.Fatal(err)
	}
}
