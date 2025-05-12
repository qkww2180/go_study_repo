package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetStudentInfo(t *testing.T) {
	id := "学生1"
	stu := GetStudentInfo(id)
	if len(stu.Name) == 0 {
		t.Fail()
	} else {
		fmt.Printf("%+v\n", stu)
	}
}

func TestGetName(t *testing.T) {
	id := "学生1"
	resp, err := http.Get("http://127.0.0.1:2345/get_name?student_id=" + id)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Println(string(bytes))
			// var stu Student
			// err:=json.Unmarshal(bytes,&stu)
			// if err!=nil{
			// 	fmt.Println(err)
			// 	t.Fail()
			// }else{
			// 	fmt.Printf("%+v\n",stu)
			// }
		}
	}
}

func TestGetAge(t *testing.T) {
	id := "学生1"
	//type Values map[string][]string
	resp, err := http.PostForm("http://127.0.0.1:2345/get_age", url.Values{"student_id": []string{id}})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Println(string(bytes))
			// var stu Student
			// err:=json.Unmarshal(bytes,&stu)
			// if err!=nil{
			// 	fmt.Println(err)
			// 	t.Fail()
			// }else{
			// 	fmt.Printf("%+v\n",stu)
			// }
		}
	}
}

func TestGetHeight(t *testing.T) {
	reader := strings.NewReader(`{"student_id":"学生1"}`)
	request, err := http.NewRequest("POST", "http://127.0.0.1:2345/get_height", reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	//当GET或POST请求需要携带Header时，只能通过http.Client发起请求，这是一种万能的请求方式
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Println(string(bytes))
			// var stu Student
			// err:=json.Unmarshal(bytes,&stu)
			// if err!=nil{
			// 	fmt.Println(err)
			// 	t.Fail()
			// }else{
			// 	fmt.Printf("%+v\n",stu)
			// }
		}
	}
}

// go test -v ./gin
