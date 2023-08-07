.PHONY: test
test:
	go test -v `go list ./... | grep -v examples` | \
    sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''


.PHONY: protoc-cln
protoc-cln:
	protoc --proto_path=clients/cln/grpc --go_out=. --go-grpc_out=. \
    --go_opt=Mnode.proto=clients/cln/grpc \
    --go_opt=Mprimitives.proto=clients/cln/grpc \
    --go-grpc_opt=Mnode.proto=clients/cln/grpc \
    --go-grpc_opt=Mprimitives.proto=clients/cln/grpc \
    clients/cln/grpc/primitives.proto clients/cln/grpc/node.proto

.PHONY: protoc-std-lightning
protoc-std-lightning:
	protoc --proto_path=pkg/standardclient/lightning/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mlightning_client.proto=pkg/standardclient/lightning \
    --go-grpc_opt=Mlightning_client.proto=pkg/standardclient/lightning \
	pkg/standardclient/lightning/proto/lightning_client.proto

.PHONY: protoc-std-bitcoin
protoc-std-bitcoin:
	protoc --proto_path=pkg/standardclient/bitcoin/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mbitcoin_client.proto=pkg/standardclient/bitcoin \
    --go-grpc_opt=Mbitcoin_client.proto=pkg/standardclient/bitcoin \
	pkg/standardclient/bitcoin/proto/bitcoin_client.proto

.PHONY: protoc-std-common
protoc-std-common:
	protoc --proto_path=pkg/standardclient/common/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mcommon_client.proto=pkg/standardclient/common \
    --go-grpc_opt=Mcommon_client.proto=pkg/standardclient/common \
	pkg/standardclient/common/proto/common_client.proto

.PHONY: protoc
protoc: protoc-cln protoc-std-lightning protoc-std-bitcoin protoc-std-common

.PHONY: generate-mocks
generate-mocks:
	go generate ./...