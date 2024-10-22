package common

import "encoding/json"

// TODO: must write struct to struct mapper instead of this function
func TypeConverter[T any](data any) (*T, error) {
	var result T
	dataJson, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}