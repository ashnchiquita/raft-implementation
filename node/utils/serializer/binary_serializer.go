package serializer

import (
	"encoding/gob"
	"os"
)

type BinarySerializer[T any] struct{}

const (
	BIN_EXTENSION = ".bin"
)

func NewBinarySerializer[T any]() *BinarySerializer[T] {
	return &BinarySerializer[T]{}
}

func (bs *BinarySerializer[T]) Serialize(data T, name string) error {
	path := GetPath(name + BIN_EXTENSION)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, FILE_PERMISSION)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (bs *BinarySerializer[T]) Deserialize(name string) (T, error) {
	var data T

	path := GetPath(name + BIN_EXTENSION)

	file, err := os.OpenFile(path, os.O_RDONLY, FILE_PERMISSION)
	if err != nil {
		return data, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)

	return data, err
}
