package serializer

import (
	"encoding/json"
	"os"
)

type JSONSerializer[T any] struct{}

const (
	JSON_EXTENSION = ".json"
)

func (js *JSONSerializer[T]) Serialize(data T, name string) error {
	path := GetPath(name + JSON_EXTENSION)

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bytes, FILE_PERMISSION)
	if err != nil {
		return err
	}

	return nil
}

func (js *JSONSerializer[T]) Deserialize(name string) (T, error) {
	var data T

	path := GetPath(name + JSON_EXTENSION)

	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(file, &data)

	return data, err
}
