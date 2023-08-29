test:
	go test -v `go list ./... | grep -v examples` | \
    sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

protoc-cln:
	protoc --proto_path=clients/cln/grpc --go_out=. --go-grpc_out=. \
    --go_opt=Mnode.proto=clients/cln/grpc \
    --go_opt=Mprimitives.proto=clients/cln/grpc \
    --go-grpc_opt=Mnode.proto=clients/cln/grpc \
    --go-grpc_opt=Mprimitives.proto=clients/cln/grpc \
    clients/cln/grpc/primitives.proto clients/cln/grpc/node.proto

protoc-std-lightning:
	protoc --proto_path=pkg/standardclient/lightning/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mlightning_client.proto=pkg/standardclient/lightning \
    --go-grpc_opt=Mlightning_client.proto=pkg/standardclient/lightning \
	pkg/standardclient/lightning/proto/lightning_client.proto

protoc-std-bitcoin:
	protoc --proto_path=pkg/standardclient/bitcoin/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mbitcoin_client.proto=pkg/standardclient/bitcoin \
    --go-grpc_opt=Mbitcoin_client.proto=pkg/standardclient/bitcoin \
	pkg/standardclient/bitcoin/proto/bitcoin_client.proto

protoc-std-common:
	protoc --proto_path=pkg/standardclient/common/proto --go_out=. --go-grpc_out=. \
    --go_opt=Mcommon_client.proto=pkg/standardclient/common \
    --go-grpc_opt=Mcommon_client.proto=pkg/standardclient/common \
	pkg/standardclient/common/proto/common_client.proto

protoc: protoc-cln protoc-std-lightning protoc-std-bitcoin protoc-std-common

build-cln-client:
	docker build -t cln-client:latest -f clients/cln/Dockerfile .

build-lnd-client:
	docker build -t lnd-client:latest -f clients/lnd/Dockerfile .

build-bitcoind-client:
	docker build -t bitcoind-client:latest -f clients/bitcoind/Dockerfile .

generate-mocks:
	go generate ./...

.PHONY: test generate-mocks protoc protoc-std-common protoc-std-bitcoin protoc-std-lightning protoc-cln build-cln-client build-lnd-client build-bitcoind-client