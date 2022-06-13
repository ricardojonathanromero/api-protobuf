binaries:
	curl -sSL \
		"https://github.com/bufbuild/buf/releases/download/v1.5.0/buf-linux-x86_64" \
		-o "/usr/local/bin/buf" && \
		chmod +x "/usr/local/bin/buf"

dependencies:
	go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

protobuf: dependencies
	buf generate

compile:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o protobuf-app cmd/main.go

release: binaries protobuf compile

kube-release:
	kubectl apply -f scripts/mongo-deployment.yml
	kubectl apply -f deployment.yml

kube-remove:
	kubectl delete -f scripts/mongo-deployment.yml
	kubectl delete -f deployment.yml
