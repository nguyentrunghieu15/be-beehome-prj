package convert

import "encoding/json"

func StructProtoToMap(req interface{}) (map[string]interface{}, error) {
	// Marshal the temporary struct to a byte array
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the byte array into a map
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
