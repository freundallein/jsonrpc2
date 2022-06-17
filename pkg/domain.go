package pkg

import (
	"context"
	"encoding/json"
	"errors"
)

const (
	ParseErrorCode     = -32700 // Invalid JSON was received by the server.
	InvalidRequestCode = -32600 // The JSON sent is not a valid Request object.
	MethodNotFoundCode = -32601 // The method does not exist / is not available.
	InvalidParamsCode  = -32602 // Invalid method parameter(s).
	InternalErrorCode  = -32603 // Internal JSON-RPC error.
)

var (
	ErrDuringParse    = errors.New("parse error")
	ErrMethodNotFound = errors.New("method not found")
	ErrInvalidRequest = errors.New("provided JSON is not a valid Request object")
	ErrInvalidParams  = errors.New("invalid method params")
)

type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Dispatcher interface {
	RegisterHandler(name string, handler interface{}) error
	DispatchMessage(ctx context.Context, rawMessage []byte) []byte
}
