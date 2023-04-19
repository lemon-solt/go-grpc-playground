# gRPC を学ぶ

# package

- github.com/golang/protobuf/jsonpb
- github.com/golang/protobuf
- google.golang.org/protobuf
- github.com/grpc-ecosystem/go-grpc-middleware
-

# vscode プラグイン

- vscode-proto3

# grpcs

- Unary RPC // 一般的なリクエスト等
- Server Streaming RPC // プッシュ通知等
- Client Streaming RPC // ファイルアップロード等
- Bidirectional Streaming RPC // リアルタイム

# docker exec

```
docker compose up -d --build
docker compose exec grpc-playground bash && cd app

```

# app

```
cd app
go run main.go
```

# relative commands

```
curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip -o protoc.zip

unzip protoc.zip


mv bin/* /usr/local/bin/ && mv include/* /usr/local/include

go install google.golang.org/protobuf/cmd/protoc-gen-go \
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

[errors]
下記エラーが出た場合
https://stackoverflow.com/questions/60578892/protoc-gen-go-grpc-program-not-found-or-is-not-executable
--go-grpc_out: protoc-gen-go-grpc: Plugin failed with status code 1.
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

-- ミドルウェア
go get -u github.com/grpc-ecosystem/go-grpc-middleware

--buffers
protoc -I. --go_out=. --go-grpc_out=. *.proto

```

# grpc document

```
package
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

cmd
protoc --doc_out=./ --doc_opt=html,sample.html file.proto
```
