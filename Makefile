protogen:
	protoc --go_out=. --go-grpc_out=. ./api/user.proto --experimental_allow_proto3_optional
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/user.proto --experimental_allow_proto3_optional

codegen:
	oapi-codegen -generate chi-server -package api api/schema.yaml > internal/generated/server.gen.go
	oapi-codegen -generate types -package api api/schema.yaml > internal/generated/models.gen.go