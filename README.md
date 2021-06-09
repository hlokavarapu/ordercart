# Service OrderCart, CustomerNotification

## Install and run services
*services run on localhost and port 50051 by default

`
git clone https://github.com/hlokavarapu/ordercart.git
cd ordercart
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
go test -v ./...
`

## Design Docs
![Alt text](designdocs/Notifier.png?raw=true "Design Doc")



