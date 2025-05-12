package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 基础数据类型枚举
const (
	TYPE_INT = iota
	TYPE_BOOL
	TYPE_FLOAT32
	TYPE_FLOAT64
	TYPE_STRING
)

// 魔数
var MAGIC = [...]byte{37, 61, 57, 23, 111} //这是个数组

// 整型转换成字节
func IntToBytes(n int) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整型
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

// 将一批参数序列化
func MarshalArguments(arguments ...any) ([]byte, error) {
	types := make([]byte, 0, len(arguments))
	lens := make([]int, 0, len(arguments))
	buf := make([]byte, 0, 256)
	buffer := bytes.NewBuffer(buf)
	for i, arg := range arguments {
		switch v := arg.(type) {
		case bool:
			err := binary.Write(buffer, binary.BigEndian, v) //v必须是固定长度的，如int64、float64，int不是固定长度，int的长度跟机器有关
			if err != nil {
				return nil, fmt.Errorf("序列化第%d个参数失败, %v, 数据类型%s", i, err, reflect.TypeOf(arg).Name())
			}
			types = append(types, TYPE_BOOL)
			lens = append(lens, 1)
		case float32:
			err := binary.Write(buffer, binary.BigEndian, v)
			if err != nil {
				return nil, fmt.Errorf("序列化第%d个参数失败, %v, 数据类型%s", i, err, reflect.TypeOf(arg).Name())
			}
			types = append(types, TYPE_FLOAT32)
			lens = append(lens, 4)
		case float64:
			err := binary.Write(buffer, binary.BigEndian, v)
			if err != nil {
				return nil, fmt.Errorf("序列化第%d个参数失败, %v, 数据类型%s", i, err, reflect.TypeOf(arg).Name())
			}
			types = append(types, TYPE_FLOAT64)
			lens = append(lens, 8)
		case int:
			err := binary.Write(buffer, binary.BigEndian, int64(v))
			if err != nil {
				return nil, fmt.Errorf("序列化第%d个参数失败, %v, 数据类型%s", i, err, reflect.TypeOf(arg).Name())
			}
			types = append(types, TYPE_INT)
			lens = append(lens, 8)
		case string:
			_, err := buffer.WriteString(v)
			if err != nil {
				return nil, fmt.Errorf("序列化第%d个参数失败, %v, 数据类型%s", i, err, reflect.TypeOf(arg).Name())
			}
			types = append(types, TYPE_STRING)
			lens = append(lens, len(v)) //断定string长度不会超过max int
		default:
			return nil, fmt.Errorf("序列化第%d个参数失败, 不支持的数据类型: %s", i, reflect.TypeOf(arg).Name())
		}
	}

	result := make([]byte, 0, len(MAGIC)+8+len(types)+8*len(types)+buffer.Len())
	writer := bytes.NewBuffer(result)
	writer.Write(MAGIC[:])               //1-魔数。[:]数组变切片
	writer.Write(IntToBytes(len(types))) //2-参数的个数
	writer.Write(types)                  //3-每一个参数的类型
	for _, l := range lens {             //4-每一个参数的长度
		writer.Write(IntToBytes(l))
	}
	writer.Write(buffer.Bytes()) //5-参数内容
	return writer.Bytes(), nil
}

// 将一批参数反序列化
func UnmarshalArguments(stream []byte) ([]any, error) {
	if len(stream) < len(MAGIC)+8 {
		return nil, fmt.Errorf("数据流长度太短: %d", len(stream))
	}
	if !bytes.Equal(MAGIC[:], stream[:len(MAGIC)]) { //1-魔数
		return nil, fmt.Errorf("魔数校验失败: %v", stream[:len(MAGIC)])
	}
	pos := len(MAGIC)
	n := BytesToInt(stream[pos : pos+8]) //2-参数的个数
	pos += 8
	if n <= 0 {
		return nil, nil
	}
	//每获得一个信息就检查一下数据长度，免得将来发生越界panic（index out of range）
	if len(stream) <= len(MAGIC)+8+n+8*n {
		return nil, fmt.Errorf("数据流长度太短: %d, 有%d个参数", len(stream), n)
	}

	types := make([]byte, 0, n) //3-每一个参数的类型
	for i := 0; i < n; i++ {
		types = append(types, stream[pos])
		pos += 1
	}

	totalLen := 0             //所有参数的总长度
	lens := make([]int, 0, n) //4-每一个参数的长度
	for i := 0; i < n; i++ {
		l := BytesToInt(stream[pos : pos+8])
		lens = append(lens, l)
		totalLen += l
		pos += 8
	}

	//每获得一个信息就检查一下数据长度，免得将来发生越界panic（index out of range）
	if len(stream[pos:]) < totalLen {
		return nil, fmt.Errorf("数据流长度太短: %d, 期待长度为: %d", len(stream), pos+totalLen)
	}

	arguments := make([]any, 0, n) //5-参数内容
	for i := 0; i < n; i++ {
		t := types[i]
		l := lens[i]
		bytesBuffer := bytes.NewBuffer(stream[pos : pos+l])

		var arg any
		switch t {
		case TYPE_BOOL:
			var x bool
			err := binary.Read(bytesBuffer, binary.BigEndian, &x)
			if err != nil {
				return nil, fmt.Errorf("反序列化第%d个参数失败, %v, 数据类型%d", i, err, t)
			}
			arg = x
		case TYPE_FLOAT32:
			var x float32
			err := binary.Read(bytesBuffer, binary.BigEndian, &x)
			if err != nil {
				return nil, fmt.Errorf("反序列化第%d个参数失败, %v, 数据类型%d", i, err, t)
			}
			arg = x
		case TYPE_FLOAT64:
			var x float64
			err := binary.Read(bytesBuffer, binary.BigEndian, &x)
			if err != nil {
				return nil, fmt.Errorf("反序列化第%d个参数失败, %v, 数据类型%d", i, err, t)
			}
			arg = x
		case TYPE_INT:
			x := BytesToInt(stream[pos : pos+l])
			arg = x
		case TYPE_STRING:
			x := string(stream[pos : pos+l])
			arg = x
		default:
			return nil, fmt.Errorf("序列化第%d个参数失败, 不支持的数据类型: %d", i, t)
		}

		pos += l
		arguments = append(arguments, arg)
	}
	return arguments, nil
}

type MySerializer struct {
}

func (s MySerializer) Marshal(object any) ([]byte, error) {
	arguments := make([]any, 0, 10)
	typ := reflect.TypeOf(object)
	v := reflect.ValueOf(object)
	for i := 0; i < v.NumField(); i++ {
		if !typ.Field(i).IsExported() { //序列化时跳过不可导出的成员
			continue
		}
		arguments = append(arguments, v.Field(i).Interface())
	}
	stream, err := MarshalArguments(arguments...) //...将切片转为不定长参数
	return stream, err
}

func (s MySerializer) Unmarshal(stream []byte, object any) error {
	typ := reflect.TypeOf(object)
	value := reflect.ValueOf(object)
	if typ.Kind() != reflect.Ptr { //因为要修改v，必须传指针
		return errors.New("必须传指针")
	}

	typ = typ.Elem() //解析指针
	value = value.Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("必须传结构体的指针")
	}

	arguments, err := UnmarshalArguments(stream)
	if err != nil {
		return err
	}
	j := 0
	for i := 0; i < typ.NumField(); i++ {
		if !typ.Field(i).IsExported() { //反序列化时跳过不可导出的成员
			continue
		}
		arg := arguments[j]
		j++
		switch typ.Field(i).Type.Kind() {
		case reflect.Int:
			if v, ok := arg.(int); ok {
				value.Field(i).SetInt(int64(v))
			} else {
				return fmt.Errorf("第%d个参数类型不一致, 期待为int", i)
			}
		case reflect.Float32:
			if v, ok := arg.(float32); ok {
				value.Field(i).SetFloat(float64(v))
			} else {
				return fmt.Errorf("第%d个参数类型不一致, 期待为float32", i)
			}
		case reflect.Float64:
			if v, ok := arg.(float64); ok {
				value.Field(i).SetFloat(v)
			} else {
				return fmt.Errorf("第%d个参数类型不一致, 期待为float64", i)
			}
		case reflect.Bool:
			if v, ok := arg.(bool); ok {
				value.Field(i).SetBool(v)
			} else {
				return fmt.Errorf("第%d个参数类型不一致, 期待为bool", i)
			}
		case reflect.String:
			if v, ok := arg.(string); ok {
				value.Field(i).SetString(v)
			} else {
				return fmt.Errorf("第%d个参数类型不一致, 期待为string", i)
			}
		default:
			return fmt.Errorf("不支持的数据类型 %s", value.Field(i).Kind().String())
		}
	}
	return nil
}
