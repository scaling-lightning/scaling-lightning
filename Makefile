.PHONY: test
test:
	go test -v `go list ./... | grep -v examples` | \
    sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

.PHONY: protoc
protoc:
	protoc --proto_path=clients/cln/grpc --go_out=. --go-grpc_out=. \
    --go_opt=Mnode.proto=clients/cln/grpc \
    --go_opt=Mprimitives.proto=clients/cln/grpc \
    --go-grpc_opt=Mnode.proto=clients/cln/grpc \
    --go-grpc_opt=Mprimitives.proto=clients/cln/grpc \
    clients/cln/grpc/primitives.proto clients/cln/grpc/node.proto
	protoc --proto_path=pkg/standardclient/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mstd_lightning_client.proto=pkg/standardclient/lightning \
    --go_opt=Mcommon.proto=pkg/standardclient/lightning \
    --go-grpc_opt=Mstd_lightning_client.proto=pkg/standardclient/lightning \
    --go-grpc_opt=Mcommon.proto=pkg/standardclient/lightning \
	pkg/standardclient/proto/common.proto pkg/standardclient/proto/std_lightning_client.proto

.PHONY: generate-mocks
generate-mocks:
	go generate ./...