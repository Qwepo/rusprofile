create:
	mkdir -p gen/openapiv2

	protoc -I proto/   \
   --go_out ./gen --go_opt paths=source_relative \
   --go-grpc_out ./gen --go-grpc_opt paths=source_relative \
   --grpc-gateway_out ./gen --grpc-gateway_opt paths=source_relative \
   --openapiv2_out ./gen/openapiv2 \
    --openapiv2_opt logtostderr=false \
	--openapiv2_opt generate_unbound_methods=true \
    proto/rusprof.proto
