module github.com/bluele/interchain-packet-router

go 1.14

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d // indirect
	github.com/bluele/interchain-simple-packet v0.0.0-20200912044232-f9777fd0845d
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200829142048-5ee4fad5010e
	github.com/gogo/protobuf v1.3.1
	github.com/grpc-ecosystem/grpc-gateway v1.14.8 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/iavl v0.14.0 // indirect
	github.com/tendermint/tendermint v0.34.0-rc3
	github.com/tendermint/tm-db v0.6.2
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
