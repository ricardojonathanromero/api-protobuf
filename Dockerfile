## build stage
FROM golang:1.17-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates make cmake curl
COPY go.mod .
RUN go mod download
COPY . .
RUN curl -sSL \
    		"https://github.com/bufbuild/buf/releases/download/v1.5.0/buf-linux-x86_64" \
    		-o "/usr/local/bin/buf" && \
    		chmod +x "/usr/local/bin/buf"
RUN go get -d google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    	go install \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
            google.golang.org/protobuf/cmd/protoc-gen-go \
            google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN buf generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o protobuf-app cmd/main.go


## deploy stage
FROM alpine as final

MAINTAINER "ricardo.jonathan.romero@gmail.com"
LABEL service="protobuf-app"
LABEL owner="ricardojonathanromero"

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
RUN adduser -S noroot
COPY --from=builder /app/protobuf-app .
COPY --from=builder /app/docs ./docs
RUN chown noroot protobuf-app && chmod +x protobuf-app
USER noroot
ENTRYPOINT ["/app/protobuf-app", "-enable_proxy=true", "-http_addr=:8080", "-grpc_addr=:8090"]
