# usage: sh scripts/protoc.sh <proto_file_name>

protoc --go_out=gRPC --go_opt=paths=source_relative \
    --go-grpc_out=gRPC --go-grpc_opt=paths=source_relative \
    node/core/$1.proto
