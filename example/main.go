package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	jsonrpc2 "github.com/freundallein/jsonrpc2/pkg"
)

// HelloParams represents input data for Hello handler
type HelloParams struct {
	Name string `json:"name"`
}

// HelloResp represents output data for Hello handler
type HelloResp struct {
	Greet string `json:"greet"`
}

// Hello is an example of JSON-RPC 2.0 handler
func Hello(ctx context.Context, params *HelloParams) (*HelloResp, error) {
	return &HelloResp{"Hello " + params.Name}, nil
}

// DummyErrorParams represents input data for DummyError handler
type DummyErrorParams struct {
	Name string `json:"name"`
}

// DummyErrorResp represents output data for DummyError handler
type DummyErrorResp struct{}

// DummyError is an example of JSON-RPC 2.0 handler, that respondes with error
func DummyError(ctx context.Context, params *DummyErrorParams) (*DummyErrorResp, error) {
	return nil, errors.New("im broken")
}

func main() {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		panic(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		panic(err)
	}
	rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"hello","params":{"name":"user"}}`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"1","result":{"greet":"Hello user"}}`)) != 0 {
		fmt.Printf("bad res '%s'\n", res)
	}
	rawMsg = []byte(`[{"jsonrpc":"2.0","id":"1","method":"hello","params":{"name":"user", "a": 1}},{"jsonrpc":"2.0","id":"2","method":"dummyError","params":{"name":"user"}}]`)
	res = disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`[{"jsonrpc":"2.0","id":"1","result":{"greet":"Hello user"}},{"jsonrpc":"2.0","id":"2","error":{"code":-32603,"message":"im broken"}}]`)) != 0 {
		fmt.Printf("Multiply messages bad res '%s'\n", res)
	}
}
