package serialization

import "github.com/bytedance/sonic"

type Json struct {
}

func (s Json) Marshal(object any) ([]byte, error) {
	stream, err := sonic.Marshal(object)
	return stream, err
}

func (s Json) Unmarshal(stream []byte, object any) error {
	err := sonic.Unmarshal(stream, object)
	return err
}
