#!/bin/sh
protoc -I. --doc_out=. --doc_opt=html,index.html *.proto
# protoc --doc_out=plugins:. --doc_opt=markdown,index.md file.proto
# protoc --doc_out=. --doc_opt=html,index.html ./app/unary-rpc-sample/file.proto