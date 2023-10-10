package utils

import (
	json2 "encoding/json"
)

func ConvertToJSON(value any) (string, error) {
	json, err := json2.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
