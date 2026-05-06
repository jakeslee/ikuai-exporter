package utils

import "encoding/json"

func ToJsonString(obj any) string {
	result, _ := json.Marshal(obj)
	return string(result)
}
