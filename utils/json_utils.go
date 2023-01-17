package utils

import (
	"encoding/json"
)

// JSONMarshall convert GO structure to bytes
func JSONMarshall(request interface{}) ([]byte, error) {

	if request == nil {
		request = map[string]string{}
	}

	data, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return data, nil
}
