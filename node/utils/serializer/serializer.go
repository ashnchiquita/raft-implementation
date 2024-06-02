package serializer

import "path/filepath"

type Serializer[T any] interface {
	Serialize(T, name string) error
	Deserialize(name string) (T, error)
}

const (
	FILE_PERMISSION = 0644
	STORAGE_PATH    = "storage"
)

func GetPath(name string) string {
	return filepath.Join(STORAGE_PATH, name)
}
