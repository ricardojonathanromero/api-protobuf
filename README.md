# api-protobuf

## considerations

execute this command to expose service via minikube
`minikube service -n backend posts-app --url`

## docker

The project could be executed using docker compose. To run the project via docker compose, run the follow
command `docker compose .`

This will create two container. The containers created are exposed below:

`mongo` - creates a mongo container that is used to persist all data. This container expose port `27017`
`app` - creates the application and expose the following ports:

    8080 - http server
    8090 - gRPC server

## TODO

- [x] Swagger documentation
- [ ] Tests (WIP)
- [ ] Documentation (WIP)