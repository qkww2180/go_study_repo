package projectprepare

// CopySlice 把切片src里的元素拷贝到dest里，返回成功拷贝的元素个数
func CopySlice[T any](dest, src []T) int {
	if len(dest) == 0 || len(src) == 0 {
		return 0
	}
	i, j := 0, 0
	for ; i < len(dest) && j < len(src); i, j = i+1, j+1 {
		dest[i] = src[j]
	}
	return i
}

var A int

type User struct {
	Name string
	age  int
}

type Ifc interface{}
