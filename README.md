# jsonrpc2

Library provides interface for implementing your own JSON-RPC 2.0 server

## Usage example
```go
package main

import (
	"bytes"
	"context"
	"fmt"
	jsonrpc2 "github.com/freundallein/jsonrpc2/pkg"
)

type HelloParams struct {
	Name string `json:"name"`
}

type HelloResp struct {
	Greet string `json:"greet"`
}

func Hello(ctx context.Context, params *HelloParams) (*HelloResp, error) {
	return &HelloResp{"Hello " + params.Name}, nil
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
    rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"hello","params":{"name":"Ivan"}}`)
    ctx := context.Background()
    res := disp.DispatchMessage(ctx, rawMsg)
    if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"1","result":{"greet":"Hello Ivan"}}`)) != 0 {
        fmt.Printf("bad res '%s'\n", res)
    }
}
```