protoc -I ./ --go_out=./ --go_opt=paths=source_relative ./*.proto
protoc -I ./ --go-grpc_out=require_unimplemented_servers=false:./ --go-grpc_opt=paths=source_relative ./*.proto
protoc-go-inject-tag -input="./*.pb.go"
