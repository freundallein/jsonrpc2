package pkg

import (
	"encoding/json"
)

func handleErrorWithCode(err error, code int64) []byte {
	response := Response{
		JsonRPC: "2.0",
		Error: &Error{
			Code:    code,
			Message: err.Error(),
		},
	}
	jsonResp, _ := json.Marshal(response)
	return jsonResp
}
