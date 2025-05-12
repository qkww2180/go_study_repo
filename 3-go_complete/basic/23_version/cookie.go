package main

import (
	"fmt"
	"net/http"
	"strings"
)

// http request header
//
// Cookie: session_id=eddycjy; value=hello-world; lang=en; lang=zh-CN
func parseCookie() {
	lines := "session_id=eddycjy; value=hello-world; lang=en; lang=zh-CN" // 存在2个lang
	cookies, _ := http.ParseCookie(lines)
	for _, cookie := range cookies {
		fmt.Printf("%s: %s\n", cookie.Name, cookie.Value)
	}
	fmt.Println(strings.Repeat("-", 50))
}

// http response header
//
// Set-Cookie: session_id=eddycjy; MaxAge=0; lang=zh-CN; Domain=.eddycjy.com
func parseSetCookie() {
	line := "session_id=eddycjy; MaxAge=0; lang=zh-CN; Domain=.eddycjy.com"
	cookie, _ := http.ParseSetCookie(line)

	fmt.Println("Name:", cookie.Name)     //session_id
	fmt.Println("Value:", cookie.Value)   //eddycjy
	fmt.Println("Domain:", cookie.Domain) //.eddycjy.com
	fmt.Println("MaxAge:", cookie.MaxAge) //0
	fmt.Println(strings.Repeat("-", 50))
}
