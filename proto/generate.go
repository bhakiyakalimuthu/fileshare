package proto

//go:generate protoc --go_out=paths=source_relative:. service.proto
//go:generate protoc --go-grpc_out=paths=source_relative:. service.proto
