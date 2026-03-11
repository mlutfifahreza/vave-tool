#!/bin/bash

# Generate Go code from protobuf definitions

echo "Generating protobuf code..."

protoc --go_out=. --go_opt=module=github.com/vave-tool/backend \
    --go-grpc_out=. --go-grpc_opt=module=github.com/vave-tool/backend \
    proto/product.proto

echo "Protobuf generation complete!"
