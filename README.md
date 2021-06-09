# Service OrderCart

## Install and run service
*service runs on localhost and port 50051 by default

`
go get https://github.com/hlokavarapu/ordercart
cd $GOPATH/src/github.com/hlokavarapu/ordercart
go run main.go
`

## Client calls:
*To make client calls use the CLI tool grpcurl to hit the OrderCart service.
For example:
`
grpcurl -plaintext -d '{"cart": [{"name": "Apple"},{"name": "Apple"},{"name": "Orange"},{"name": "Apple"}]}' localhost:50051 ordercart.OrderCart/GetOrderCost
`

## Run tests
`
go test ./...
`


