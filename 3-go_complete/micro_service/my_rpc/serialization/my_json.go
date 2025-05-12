package serialization

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MyJson struct {
}

// 由于json字符串里存在{}[]等嵌套情况，直接按,分隔是不合适的
func (s MyJson) SplitJson(json string) []string {
	rect := make([]string, 0, 10)
	stack := list.New() //list是双端队列，用它来模拟栈
	beginIndex := 0
	for i, r := range json {
		if r == rune('{') || r == rune('[') {
			stack.PushBack(struct{}{}) //我们不关心栈里是什么，只关心栈里有没有元素
		} else if r == rune('}') || r == rune(']') {
			ele := stack.Back()
			if ele != nil {
				stack.Remove(ele) //删除栈顶元素
			}
		} else if r == rune(',') {
			if stack.Len() == 0 { //栈为空时才可以按,分隔
				rect = append(rect, json[beginIndex:i])
				beginIndex = i + 1
			}
		}
	}
	rect = append(rect, json[beginIndex:])
	return rect
}

func (inst MyJson) Marshal(v any) ([]byte, error) {
	value := reflect.ValueOf(v)
	typ := value.Type() //跟typ := reflect.TypeOf(v)等价
	if typ.Kind() == reflect.Ptr {
		if value.IsNil() { //如果指向nil，直接输出null
			return []byte("null"), nil
		} else { //如果传的是指针类型，先解析指针
			typ = typ.Elem()
			value = value.Elem()
		}
	}
	bf := bytes.Buffer{} //存放序列化结果
	switch typ.Kind() {
	case reflect.String:
		return []byte(fmt.Sprintf("\"%s\"", value.String())), nil //取得reflect.Value对应的原始数据的值
	case reflect.Bool:
		return []byte(fmt.Sprintf("%t", value.Bool())), nil
	case reflect.Float32,
		reflect.Float64:
		return []byte(fmt.Sprintf("%f", value.Float())), nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return []byte(fmt.Sprintf("%v", value.Interface())), nil
	case reflect.Slice:
		if value.IsNil() {
			return []byte("null"), nil
		}
		bf.WriteByte('[')
		if value.Len() > 0 {
			for i := 0; i < value.Len(); i++ { //取得slice的长度
				if bs, err := inst.Marshal(value.Index(i).Interface()); err != nil { //对slice的第i个元素进行序列化。递归
					return nil, err
				} else {
					bf.Write(bs)
					bf.WriteByte(',')
				}
			}
			bf.Truncate(len(bf.Bytes()) - 1) //删除最后一个逗号
		}
		bf.WriteByte(']')
		return bf.Bytes(), nil
	case reflect.Map:
		if value.IsNil() {
			return []byte("null"), nil
		}
		bf.WriteByte('{')
		if value.Len() > 0 {
			for _, key := range value.MapKeys() {
				if keyBs, err := inst.Marshal(key.Interface()); err != nil {
					return nil, err
				} else {
					bf.Write(keyBs)
					bf.WriteByte(':')
					v := value.MapIndex(key)
					if vBs, err := inst.Marshal(v.Interface()); err != nil {
						return nil, err
					} else {
						bf.Write(vBs)
						bf.WriteByte(',')
					}
				}
			}
			bf.Truncate(len(bf.Bytes()) - 1) //删除最后一个逗号
		}
		bf.WriteByte('}')
		return bf.Bytes(), nil
	case reflect.Struct:
		bf.WriteByte('{')
		if value.NumField() > 0 {
			for i := 0; i < value.NumField(); i++ {
				fieldValue := value.Field(i)
				fieldType := typ.Field(i)
				name := fieldType.Name //如果没有json Tag，默认使用成员变量的名称
				if len(fieldType.Tag.Get("json")) > 0 {
					name = fieldType.Tag.Get("json")
				}
				bf.WriteString("\"")
				bf.WriteString(name)
				bf.WriteString("\"")
				bf.WriteString(":")
				if bs, err := inst.Marshal(fieldValue.Interface()); err != nil { //对value递归调用Marshal序列化
					return nil, err
				} else {
					bf.Write(bs)
				}
				bf.WriteString(",")
			}
			bf.Truncate(len(bf.Bytes()) - 1) //删除最后一个逗号
		}
		bf.WriteByte('}')
		return bf.Bytes(), nil
	default:
		return []byte(fmt.Sprintf("\"暂不支持该数据类型:%s\"", typ.Kind().String())), nil
	}
}

func (inst MyJson) Unmarshal(data []byte, v any) error {
	s := string(data)
	//去除前后的连续空格
	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")
	if len(s) == 0 {
		return nil
	}
	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)
	if typ.Kind() != reflect.Ptr { //因为要修改v，必须传指针
		return errors.New("must pass pointer parameter")
	}

	typ = typ.Elem() //解析指针
	value = value.Elem()

	switch typ.Kind() {
	case reflect.String:
		if s[0] == '"' && s[len(s)-1] == '"' {
			value.SetString(s[1 : len(s)-1]) //去除前后的""
		} else {
			return fmt.Errorf("invalid json part: %s", s)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err == nil {
			value.SetBool(b)
		} else {
			return err
		}
	case reflect.Float32,
		reflect.Float64:
		if f, err := strconv.ParseFloat(s, 64); err != nil {
			return err
		} else {
			value.SetFloat(f) //通过reflect.Value修改原始数据的值
		}
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if i, err := strconv.ParseInt(s, 10, 64); err != nil {
			return err
		} else {
			value.SetInt(i) //有符号整型通过SetInt
		}
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if i, err := strconv.ParseUint(s, 10, 64); err != nil {
			return err
		} else {
			value.SetUint(i) //无符号整型需要通过SetUint
		}
	case reflect.Slice:
		if s[0] == '[' && s[len(s)-1] == ']' {
			arr := inst.SplitJson(s[1 : len(s)-1]) //去除前后的[]
			if len(arr) > 0 {
				slice := reflect.ValueOf(v).Elem()                    //别忘了，v是指针
				slice.Set(reflect.MakeSlice(typ, len(arr), len(arr))) //通过反射创建slice
				for i := 0; i < len(arr); i++ {
					eleValue := slice.Index(i)
					eleType := eleValue.Type()
					if eleType.Kind() != reflect.Ptr {
						eleValue = eleValue.Addr()
					}
					if err := inst.Unmarshal([]byte(arr[i]), eleValue.Interface()); err != nil {
						return err
					}
				}
			}
		} else if s != "null" {
			return fmt.Errorf("invalid json part: %s", s)
		}
	case reflect.Map:
		if s[0] == '{' && s[len(s)-1] == '}' {
			arr := inst.SplitJson(s[1 : len(s)-1]) //去除前后的{}
			if len(arr) > 0 {
				mapValue := reflect.ValueOf(v).Elem()                //别忘了，v是指针
				mapValue.Set(reflect.MakeMapWithSize(typ, len(arr))) //通过反射创建map
				kType := typ.Key()                                   //获取map的key的Type
				vType := typ.Elem()                                  //获取map的value的Type
				for i := 0; i < len(arr); i++ {
					brr := strings.Split(arr[i], ":")
					if len(brr) != 2 {
						return fmt.Errorf("invalid json part: %s", arr[i])
					}

					kValue := reflect.New(kType) //根据Type创建指针型的Value
					if err := inst.Unmarshal([]byte(brr[0]), kValue.Interface()); err != nil {
						return err
					}
					vValue := reflect.New(vType) //根据Type创建指针型的Value
					if err := inst.Unmarshal([]byte(brr[1]), vValue.Interface()); err != nil {
						return err
					}
					mapValue.SetMapIndex(kValue.Elem(), vValue.Elem()) //往map里面赋值
				}
			}
		} else if s != "null" {
			return fmt.Errorf("invalid json part: %s", s)
		}
	case reflect.Struct:
		if s[0] == '{' && s[len(s)-1] == '}' {
			arr := inst.SplitJson(s[1 : len(s)-1])
			if len(arr) > 0 {
				fieldCount := typ.NumField()
				//建立json tag到FieldName的映射关系
				tag2Field := make(map[string]string, fieldCount)
				for i := 0; i < fieldCount; i++ {
					fieldType := typ.Field(i)
					name := fieldType.Name
					if len(fieldType.Tag.Get("json")) > 0 {
						name = fieldType.Tag.Get("json")
					}
					tag2Field[name] = fieldType.Name
				}

				for _, ele := range arr {
					brr := strings.SplitN(ele, ":", 2) //json的value里可能存在嵌套，所以用:分隔时限定个数为2
					if len(brr) == 2 {
						tag := strings.Trim(brr[0], " ")
						if tag[0] == '"' && tag[len(tag)-1] == '"' { //json的key肯定是带""的
							tag = tag[1 : len(tag)-1]                        //去除json key前后的""
							if fieldName, exists := tag2Field[tag]; exists { //根据json key(即json tag)找到对应的FieldName
								fieldValue := value.FieldByName(fieldName)
								fieldType := fieldValue.Type()
								if fieldType.Kind() != reflect.Ptr {
									//如果内嵌不是指针，则声明时已经用0值初始化了，此处只需要根据json改写它的值
									fieldValue = fieldValue.Addr()                                                 //确保fieldValue指向指针类型，因为接下来要把fieldValue传给Unmarshal
									if err := inst.Unmarshal([]byte(brr[1]), fieldValue.Interface()); err != nil { //递归调用Unmarshal，给fieldValue的底层数据赋值
										return err
									}
								} else {
									//如果内嵌的是指针，则需要通过New()创建一个实例(申请内存空间)。不能给New()传指针型的Type，所以调一下Elem()
									newValue := reflect.New(fieldType.Elem())                                    //newValue代表的是指针
									if err := inst.Unmarshal([]byte(brr[1]), newValue.Interface()); err != nil { //递归调用Unmarshal，给fieldValue的底层数据赋值
										return err
									}
									value.FieldByName(fieldName).Set(newValue) //把newValue赋给value的Field
								}

							} else {
								fmt.Printf("字段%s找不到\n", tag)
							}
						} else {
							return fmt.Errorf("invalid json part: %s", tag)
						}
					} else {
						return fmt.Errorf("invalid json part: %s", ele)
					}
				}
			}
		} else if s != "null" {
			return fmt.Errorf("invalid json part: %s", s)
		}
	default:
		fmt.Printf("暂不支持类型:%s\n", typ.Kind().String())
	}
	return nil
}
