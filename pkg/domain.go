package pkg

import (
	"context"
	"encoding/json"
	"errors"
)

const (
	// ParseErrorCode for errors when invalid JSON was received by the server.
	ParseErrorCode     = -32700
	// InvalidRequestCode for errors when provided is not a valid Request object.
	InvalidRequestCode = -32600
	// MethodNotFoundCode for errors when method does not exist
	MethodNotFoundCode = -32601
	// InvalidParamsCode for errors when received invalid method parameter(s).
	InvalidParamsCode  = -32602
	// InternalErrorCode for errors when server has internal JSON-RPC error
	InternalErrorCode  = -32603
)

var (
	// ErrDuringParse can be raised during json body parsing
	ErrDuringParse    = errors.New("parse error")
	// ErrMethodNotFound raised if no JSON-RPC handler found
	ErrMethodNotFound = errors.New("method not found")
	// ErrInvalidRequest raised for invalid JSON-RPC requests
	ErrInvalidRequest = errors.New("provided JSON is not a valid Request object")
	// ErrInvalidParams raised for invalid JSON-RPC method params
	ErrInvalidParams  = errors.New("invalid method params")
)

// Request represents JSON-RPC Request object
type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// Response represents JSON-RPC Response object
type Response struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error used for proper JSON-RPC error representation
type Error struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Dispatcher describes library user's interface
type Dispatcher interface {
	// RegisterHandler allows adding a new handler for some action
	RegisterHandler(name string, handler interface{}) error
	// DispatchMessage serves raw JSON-RPC message, executes required action and returns result
	DispatchMessage(ctx context.Context, rawMessage []byte) []byte
}
