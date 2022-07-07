.PHONY: kube-release kube-remove
kube-release:
	docker build -t app .
	kubectl apply -f scripts/mongo-deployment.yml
	kubectl apply -f deployment.yml

kube-remove:
	kubectl delete -f scripts/mongo-deployment.yml
	kubectl delete -f deployment.yml

.PHONY: tests
tests:
	go test ./...

.PHONY: dependencies
dependencies:
	curl -sSL \
        		"https://github.com/bufbuild/buf/releases/download/v1.5.0/buf-linux-x86_64" \
        		-o "/usr/local/bin/buf" && \
        		chmod +x "/usr/local/bin/buf"
	go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
		go install \
			github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
			github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
			google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/jstemmer/go-junit-report/v2@latest \
			github.com/axw/gocov/gocov@latest \
			github.com/AlekSi/gocov-xml@latest
	mkdir coverage/reports
