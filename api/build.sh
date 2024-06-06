# protoc -I . \
#     --go_out . --go_opt paths=source_relative \
#     --go-grpc_out . --go-grpc_opt paths=source_relative \
#     --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
#     --openapiv2_out . \
#     --openapiv2_opt use_go_templates=true \
#     ./api/auth-api/auth_api.proto

# protoc -I . \
#     --go_out . --go_opt paths=source_relative \
#     --go-grpc_out . --go-grpc_opt paths=source_relative \
#     --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
#     --openapiv2_out . \
#     --openapiv2_opt use_go_templates=true \
#     ./api/user-api/user_api.proto

protoc -I . \
    --go_out . --go_opt paths=source_relative \
    ./api/pro-api/message.proto

protoc -I . \
    --go_out . --go_opt paths=source_relative \
    --go-grpc_out . --go-grpc_opt paths=source_relative \
    --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
    --openapiv2_out . \
    --openapiv2_opt use_go_templates=true,allow_merge=true,merge_file_name=./api/pro-api/pro-api \
    ./api/pro-api/*.proto