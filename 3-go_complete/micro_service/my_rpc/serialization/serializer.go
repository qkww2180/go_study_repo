package serialization

type Serializer interface {
	Marshal(object any) ([]byte, error)
	Unmarshal(stream []byte, object any) error
}
