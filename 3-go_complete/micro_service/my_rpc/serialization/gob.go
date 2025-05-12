package serialization

import (
	"bytes"
	"encoding/gob"
)

type Gob struct {
}

func (s Gob) Marshal(object any) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(object)
	return buffer.Bytes(), err
}

func (s Gob) Unmarshal(stream []byte, object any) error {
	buffer := bytes.NewBuffer(stream)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(object)
	return err
}
