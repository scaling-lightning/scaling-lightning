.PHONY: test
test:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

.PHONY: protoc
protoc:
	protoc --proto_path=clients/cln/grpc --go_out=. --go-grpc_out=. \
    --go_opt=Mnode.proto=clients/cln/grpc \
    --go_opt=Mprimitives.proto=clients/cln/grpc \
    --go-grpc_opt=Mnode.proto=clients/cln/grpc \
    --go-grpc_opt=Mprimitives.proto=clients/cln/grpc \
    clients/cln/grpc/primitives.proto clients/cln/grpc/node.proto

.PHONY: generate-mocks
generate-mocks:
	go generate ./...