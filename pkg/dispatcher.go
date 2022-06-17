package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// NewDispatcher is a constructor for a new Dispatcher instance
func NewDispatcher() Dispatcher {
	return &JSONRPCv2Dispatcher{
		handlers: map[string]interface{}{},
	}
}

// JSONRPCv2Dispatcher implementation of JSON-RPC v2.0 server
type JSONRPCv2Dispatcher struct {
	handlers map[string]interface{}
}

// RegisterHandler allows adding a new handler for some action
func (d *JSONRPCv2Dispatcher) RegisterHandler(name string, handler interface{}) error {
	if _, ok := d.handlers[name]; ok {
		return errors.New(fmt.Sprintf("handler with name %s is already registered", name))
	}
	err := d.validateHandler(name, handler)
	if err != nil {
		return err
	}
	d.handlers[name] = handler
	return nil
}

func (d *JSONRPCv2Dispatcher) validateHandler(name string, handler interface{}) error {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return errors.New(fmt.Sprintf("%s expected function, but got: %v", name, handlerType.Kind()))
	}
	// Input args validation
	if handlerType.NumIn() != 2 {
		return errors.New(fmt.Sprintf("should be 2 input params in handler %s", name))
	}
	ctxType := handlerType.In(0)
	ctxInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !ctxType.Implements(ctxInterface) {
		return errors.New(fmt.Sprintf("%s expected context.Context as first input param, but got: %v", name, ctxType))
	}
	inputParamsType := handlerType.In(1)
	if inputParamsType.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("%s expected pointer to struct as second input param, but got: %v", name, inputParamsType.Kind()))
	}
	if inputParamsType.Elem().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("%s expected pointer to struct as second input param, but got ptr to: %v", name, inputParamsType.Elem().Kind()))
	}
	// Output args validation
	if handlerType.NumOut() != 2 {
		return errors.New(fmt.Sprintf("should be 2 output params in handler %s", name))
	}
	outputParamsType := handlerType.Out(0)
	if outputParamsType.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("%s expected pointer to struct as first output param, but got: %v", name, outputParamsType.Kind()))
	}
	if outputParamsType.Elem().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("%s expected pointer to struct as first output param, but got ptr to: %v", name, outputParamsType.Elem().Kind()))
	}
	errType := handlerType.Out(1)
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !errType.Implements(errInterface) {
		return errors.New(fmt.Sprintf("%s expected error as a second output param, but got: %v", name, errType))
	}
	return nil
}

// DispatchMessage serves raw JSON-RPC message, executes required action and returns result
func (d *JSONRPCv2Dispatcher) DispatchMessage(ctx context.Context, rawMessage []byte) []byte {
	requests, err := d.parseMessage(rawMessage)
	if err != nil {
		return handleErrorWithCode(err, ParseErrorCode)
	}
	responses := []*Response{}
	for _, request := range requests {
		response := &Response{
			JsonRPC: "2.0",
			ID:      request.ID,
		}
		resp, err := d.handleRequest(ctx, request)
		if err != nil {
			switch err {
			case ErrInvalidRequest:
				response.Error = &Error{
					Code:    InvalidRequestCode,
					Message: err.Error(),
				}
			case ErrMethodNotFound:
				response.Error = &Error{
					Code:    MethodNotFoundCode,
					Message: err.Error(),
				}
			case ErrInvalidParams:
				response.Error = &Error{
					Code:    InvalidParamsCode,
					Message: err.Error(),
				}
			default:
				response.Error = &Error{
					Code:    InternalErrorCode,
					Message: err.Error(),
				}
			}
		} else {
			response.Result = resp
		}
		responses = append(responses, response)
	}
	return d.marshalResponses(responses)
}

func (d *JSONRPCv2Dispatcher) parseMessage(message []byte) ([]*Request, error) {
	requests := []*Request{}
	if message[0] == []byte("{")[0] {
		message = append([]byte("["), message...)
		message = append(message, []byte("]")[0])
	}
	var rawRequests []*json.RawMessage
	err := json.Unmarshal(message, &rawRequests)
	if err != nil {
		return nil, ErrDuringParse
	}
	for _, rawRequest := range rawRequests {
		request := Request{}
		// skip json error, will handle empty request as invalid
		_ = json.Unmarshal([]byte(*rawRequest), &request)
		requests = append(requests, &request)
	}
	return requests, nil
}

func (d *JSONRPCv2Dispatcher) marshalResponses(responses []*Response) []byte {
	if len(responses) == 1 {
		response, err := json.Marshal(responses[0])
		if err != nil {
			return handleErrorWithCode(err, InternalErrorCode)
		}
		return response
	}
	response, err := json.Marshal(responses)
	if err != nil {
		return handleErrorWithCode(err, InternalErrorCode)
	}
	return response
}

func (d *JSONRPCv2Dispatcher) handleRequest(ctx context.Context, req *Request) (interface{}, error) {
	if req.Method == "" || req.JsonRPC != "2.0" {
		return nil, ErrInvalidRequest
	}
	handler, exists := d.handlers[req.Method]
	if !exists {
		return nil, ErrMethodNotFound
	}
	v := reflect.ValueOf(handler)
	handlerType := reflect.TypeOf(handler)
	inputParamsPtr := reflect.New(handlerType.In(1))
	err := json.Unmarshal(req.Params, inputParamsPtr.Interface())
	if err != nil {
		return nil, ErrInvalidParams
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		inputParamsPtr.Elem(),
	}
	vResponse := v.Call(args)
	result := vResponse[0].Interface()
	vErr := vResponse[1].Interface()
	if vErr != nil {
		return nil, vErr.(error)
	}
	return result, nil
}
