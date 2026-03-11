#!/bin/bash

# Generate Go code from protobuf definitions

echo "Generating protobuf code..."

protoc --go_out=internal/grpc --go_opt=paths=source_relative \
    --go-grpc_out=internal/grpc --go-grpc_opt=paths=source_relative \
    proto/product.proto

echo "Protobuf generation complete!"
