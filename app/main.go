package main

import (
	unaryrpcsample "app/unary-rpc-sample"
)

func main() {
	// marhsal proto
	// grpcmarshalsample.GrpcMarshalSample()

	// unery RPC call
	unaryrpcsample.SwichClientOrServer()
}
