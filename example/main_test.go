package main

import (
	"bytes"
	"context"
	jsonrpc2 "github.com/freundallein/jsonrpc2/pkg"
	"testing"
)

func TestDispatchMessage(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"hello","params":{"name":"user"}}`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"1","result":{"greet":"Hello user"}}`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
	rawMsg = []byte(`[{"jsonrpc":"2.0","id":"1","method":"hello","params":{"name":"user", "a": 1}},{"jsonrpc":"2.0","id":"2","method":"dummyError","params":{"name":"user"}}]`)
	res = disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`[{"jsonrpc":"2.0","id":"1","result":{"greet":"Hello user"}},{"jsonrpc":"2.0","id":"2","error":{"code":-32603,"message":"im broken"}}]`)) != 0 {
		t.Errorf("Multiply messages bad result '%s'\n", res)
	}
}

func TestParseError(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"he`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"","error":{"code":-32700,"message":"parse error"}}`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
}

func TestInvalidRequest(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`{}`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"","error":{"code":-32600,"message":"provided JSON is not a valid Request object"}}`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
}

func TestInvalidRequestBatch(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`[1,2,3]`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`[{"jsonrpc":"2.0","id":"","error":{"code":-32600,"message":"provided JSON is not a valid Request object"}},{"jsonrpc":"2.0","id":"","error":{"code":-32600,"message":"provided JSON is not a valid Request object"}},{"jsonrpc":"2.0","id":"","error":{"code":-32600,"message":"provided JSON is not a valid Request object"}}]`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
}

func TestMethodNotFound(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"unknownMethod"}`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"1","error":{"code":-32601,"message":"method not found"}}`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
}

func TestInvalidParams(t *testing.T) {
	disp := jsonrpc2.NewDispatcher()
	err := disp.RegisterHandler("hello", Hello)
	if err != nil {
		t.Error(err)
	}
	err = disp.RegisterHandler("dummyError", DummyError)
	if err != nil {
		t.Error(err)
	}
	rawMsg := []byte(`{"jsonrpc": "2.0","id":"1","method":"hello","params":1}`)
	ctx := context.Background()
	res := disp.DispatchMessage(ctx, rawMsg)
	if bytes.Compare(res, []byte(`{"jsonrpc":"2.0","id":"1","error":{"code":-32602,"message":"invalid method params"}}`)) != 0 {
		t.Errorf("bad result '%s'\n", res)
	}
}
