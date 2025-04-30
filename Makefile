protogen:
	protoc --go_out=. --go-grpc_out=. ./api/user.proto
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/user.proto