package JSONSchemaReader

import (
	"os"
)

type Reader struct {
}

func New() *Reader {
	return &Reader{}
}

func (r *Reader) ReadMappings(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
