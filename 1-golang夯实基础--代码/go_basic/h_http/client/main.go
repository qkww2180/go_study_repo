package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	myhttp "dqq/go/basic/h_http"
)

// 详细看一下http协议
func HttpObservation() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/obs?name=dqq"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		fmt.Printf("response proto: %s\n", resp.Proto)
		if major, minor, ok := http.ParseHTTPVersion(resp.Proto); ok {
			fmt.Printf("h_http major version %d, h_http minor version %d\n", major, minor)
		}

		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response header")
		for key, values := range resp.Header {
			fmt.Printf("%s: %v\n", key, values)
			if key == "Date" {
				if tm, err := http.ParseTime(values[0]); err == nil {
					fmt.Printf("server time %s\n", tm.Format("2006-01-02 15:04:05"))
				}
			}
		}
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n\n")
	}
}

func Get() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/get?" + myhttp.EncodeUrlParams(map[string]string{"name": "朝阳 Zhang", "age": "18"})); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		// j_io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		if body, err := io.ReadAll(resp.Body); err == nil {
			fmt.Print(string(body))
		}
		os.Stdout.WriteString("\n\n")
	}
}

// 流式传输海量数据
func HugeBody() {
	fmt.Println(strings.Repeat("*", 30) + "GET HUGE BODY" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/stream"); err != nil {
		panic(err)
	} else {
		headerKey := http.CanonicalHeaderKey("content-length") // 正规化之后是Content-Length
		if ls, exists := resp.Header[headerKey]; exists {
			if l, err := strconv.Atoi(ls[0]); err == nil {
				fmt.Printf("Content-Length=%d\n", l)
				total := 0
				reader := bufio.NewReader(resp.Body)
				for total < l {
					if bs, err := reader.ReadBytes('\n'); err == nil {
						total += len(bs)
						fmt.Printf("进度 %.2f%%, 内容 %s", 100*float64(total)/float64(l), string(bs)) // bs末尾包含了\n
					} else {
						if err == io.EOF {
							if len(bs) > 0 { // 即使读到末尾了，本次read也可能读出了内容
								total += len(bs)
								fmt.Printf("进度 %.2f%%, 内容 %s", 100*float64(total)/float64(l), string(bs))
							}
						} else {
							fmt.Printf("read response body error: %s\n", err)
						}
						break
					}
					if total >= l/2 {
						resp.Body.Close()
						break
					}
				}
			}
		}
		resp.Body.Close()
	}
}

func Post() {
	fmt.Println(strings.Repeat("*", 30) + "POST" + strings.Repeat("*", 30))
	//Content-Type为text/plain，表示一个朴素的字符串
	if resp, err := http.Post("http://127.0.0.1:5678/post", "text/plain", strings.NewReader("Hello Server")); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}

	//Content-Type为application/json，表示一个json字符串
	bs, _ := json.Marshal(map[string]string{"name": "朝阳 Zhang", "age": "18"})
	if resp, err := http.Post("http://127.0.0.1:5678/post", "application/json", bytes.NewReader(bs)); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}

	// PostForm()会自动把请求头的Content-Type设为application/x-www-form-urlencoded，并把url.Values转为URL-encoded参数格式放到请求体里
	if resp, err := http.PostForm("http://127.0.0.1:5678/post", url.Values{"name": []string{"朝阳 Zhang"}, "age": []string{"18"}}); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}
}

// head通常用于检测某个网页是否存在
func Head() {
	fmt.Println(strings.Repeat("*", 30) + "HEAD" + strings.Repeat("*", 30))
	//HEAD类似于GET，但HEAD方法只能取到响应头，不能取到响应体
	if resp, err := http.Head("http://127.0.0.1:5678/get"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}

	if resp, err := http.Head("https://www.baidu.com"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}

	if resp, err := http.Head("http://127.0.0.1:5678/strange"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}
}

func Cookie() {
	fmt.Println(strings.Repeat("*", 30) + "COOKIE" + strings.Repeat("*", 30))
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:5678/cookie", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (x64)") //伪造User-Agent，爬虫经常这么干
	request.Header.Add("user-role", "vip")                //header的key和value可以随意设置
	//  可添加多个Cookie
	request.AddCookie(
		&http.Cookie{
			Name:   "auth",
			Value:  "pass",
			Domain: "localhost",
			Path:   "/",
		},
	)
	//所有的cookie都会放到一个http request header中。Cookie: [auth=pass; money=100dollar]
	request.AddCookie(&http.Cookie{
		Name:  "money",
		Value: "100",
	})
	request.AddCookie(&http.Cookie{
		Name:  "money", //cookie的Name允许有重复，不会覆盖
		Value: "800",
	})
	//设置请求超时
	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}
	//发起请求
	if resp, err := client.Do(request); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response header:")
		// 其实可以直接通过resp.Cookies()获得*h_http.Cookie，没必要自己解析
		if values, exists := resp.Header["Set-Cookie"]; exists {
			fmt.Println(values[0])
			cookie, _ := http.ParseSetCookie(values[0])
			fmt.Println("Name:", cookie.Name)
			fmt.Println("Value:", cookie.Value)
			fmt.Println("Domain:", cookie.Domain)
			fmt.Println("MaxAge:", cookie.MaxAge)
			fmt.Println(strings.Repeat("-", 50))
		}
		os.Stdout.WriteString("\n\n")
	}
}

func Restful() {
	fmt.Println(strings.Repeat("*", 30) + "GET RESTful" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/get/zcy/18"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		// j_io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		if body, err := io.ReadAll(resp.Body); err == nil {
			fmt.Print(string(body))
		}
		os.Stdout.WriteString("\n\n")
	}
}

func Student() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/student"); err != nil { // 直接在浏览器里访问http://127.0.0.1:5678/student
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n\n")
	}
}

func main() {
	HttpObservation()
	Get()
	HugeBody()
	Post()
	Head()
	Cookie()
	Restful()
	Student()
}

// go run ./h_http/client
