package orm

import (
	"fmt"
	"testing"
)

func TestPrintFieldInfo(t *testing.T) {
	PrintFieldInfo(Student{})
}

func TestPrintFieldValue(t *testing.T) {
	s := Student{City: "大连"}
	PrintFieldValue(&s) //得传指针，否则无法修改City的值
	fmt.Println("修改后：" + s.City)
}

func TestExecMethod(t *testing.T) {
	ExecMethod(Student{})
}

// go test -v ./g_database/orm/ -run=^TestPrintFieldInfo$ -count=1
// go test -v ./g_database/orm/ -run=^TestPrintFieldValue$ -count=1
// go test -v ./g_database/orm/ -run=^TestExecMethod$ -count=1
