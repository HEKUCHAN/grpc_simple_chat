```
# サーバー用
protoc --proto_path=protos \
  --go_out=paths=source_relative:apps/server/proto \
  --go-grpc_out=paths=source_relative:apps/server/proto \
  protos/chat.proto

# Go クライアント用
protoc --proto_path=protos \
  --go_out=paths=source_relative:apps/golang_client/proto \
  --go-grpc_out=paths=source_relative:apps/golang_client/proto \
  protos/chat.proto

# Python  クライアント用
python -m grpc_tools.protoc -Iprotos \
  --python_out=apps/python_client/proto \
  --grpc_python_out=apps/python_client/proto \
  protos/chat.proto
```
