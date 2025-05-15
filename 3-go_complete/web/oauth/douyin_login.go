// 抖音Web授权文档：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/sdk/web-app/web/permission
package main

import (
	"dqq/web/util"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

var (
	clientKey    string
	clientSecret string
)

func Init() {
	configReader := util.CreateConfigReader("auth.yaml") //用viper读取配置文件
	clientKey = configReader.GetString("douyin.ClientKey")
	clientSecret = configReader.GetString("douyin.ClientSecret")
	if len(clientKey) == 0 || len(clientSecret) == 0 {
		panic("either clientKey or clientSecret is empty")
	}
}

type tokenRequest struct {
	Secret string `json:"client_secret"`
	Code   string `json:"code"`
	Grant  string `json:"grant_type"`
	Key    string `json:"client_key"`
}

type douyinData struct {
	AccessToken string `json:"access_token"`
	OpenId      string `json:"open_id"`
	Description string `json:"description"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
}

type douyinResponse struct {
	Data    douyinData `json:"z_data"`
	Message string     `json:"message"`
}

// step3. 拿Code去调用抖音，获取access_token和open_id。调用文档https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-permission/get-access-token/
func getAccessToken(code string) (string, string, error) {
	request := &tokenRequest{
		Secret: clientSecret,
		Code:   code,
		Grant:  "authorization_code",
		Key:    clientKey,
	}
	reqStr, err := sonic.MarshalString(request)
	if err != nil {
		return "", "", err
	}
	reader := strings.NewReader(reqStr)

	if req, err := http.NewRequest(http.MethodPost, "https://open.douyin.com/oauth/access_token/", reader); err != nil {
		return "", "", err
	} else {
		req.Header.Add("Content-Type", "application/json")
		client := http.Client{}
		if resp, err := client.Do(req); err != nil {
			return "", "", err
		} else {
			defer resp.Body.Close()
			if bs, err := io.ReadAll(resp.Body); err != nil {
				return "", "", err
			} else {
				var token douyinResponse
				if err := sonic.Unmarshal(bs, &token); err != nil {
					return "", "", err
				} else {
					if token.Message == "success" {
						return token.Data.AccessToken, token.Data.OpenId, nil
					} else {
						return "", "", errors.New(token.Data.Description)
					}
				}
			}
		}
	}
}

// step4. 根据access_token和open_id获取用户的公开信息。调用文档https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/get-account-open-info
func getUserInfo(accessToken, openId string) (string, string, error) {
	if resp, err := http.PostForm("https://open.douyin.com/oauth/userinfo/", url.Values{"open_id": []string{openId}, "access_token": []string{accessToken}}); err != nil {
		return "", "", err
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		if bs, err := io.ReadAll(resp.Body); err != nil {
			return "", "", err
		} else {
			var response douyinResponse
			if err := sonic.Unmarshal(bs, &response); err != nil {
				return "", "", err
			} else {
				if response.Message == "success" {
					return response.Data.Nickname, response.Data.Avatar, nil
				} else {
					return "", "", errors.New(response.Data.Description)
				}
			}
		}
	}
}

// Step1. 跳转到抖音登录页（扫码登录或手机验证码登录）。请求链接“https://open.douyin.com/platform/oauth/connect?client_key=CLIENT_KEY&response_type=code&scope=user_info&redirect_uri=REDIRECT_URI”，注册抖音开发者账号并创建Web应用后会拿到CLIENT_KEY，REDIRECT_URI是我们准备接收（抖音返回的）code的url，在本例中就是https://localhost:5678/douyin_login
//
// scope: user_info - 用户公开信息，mobile_alert - 用户手机号，fans.check - 粉丝判断，op.business.status - 用户经营身份。官方文档参考 https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/openapi/account-management/get-account-open-info
//
// 对应的URL是https://localhost:5678/douyin_login
func douyinLogin(ctx *gin.Context) {
	//step2. 从URL参数中获取code
	code := ctx.Query("code")
	if len(code) == 0 {
		log.Printf("douyin return empty code")
		return
	}
	//拿code去调用抖音，获取access_token和open_id。code只有10分钟有效期，且只能使用一次
	accessToken, openId, err := getAccessToken(code)
	if err != nil {
		log.Printf("get access token failed: %s", err.Error())
		return
	}
	//拿access_token和open_id获取用户信息
	nickName, avatar, err := getUserInfo(accessToken, openId)
	if err != nil {
		log.Printf("get user info failed: %s", err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "logged_in.html", gin.H{"nick_name": nickName, "avatar": avatar})
}

func main() {
	Init() //初始化

	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "localhost:5678",
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	// gin.SetMode(gin.ReleaseMode) //GIN线上发布模式
	gin.DefaultWriter = io.Discard //禁止GIN的输出
	engine := gin.Default()

	engine.Use(secureFunc)
	engine.LoadHTMLFiles("web/oauth/static/login.html", "web/oauth/static/logged_in.html")
	engine.GET("/", func(ctx *gin.Context) { //请求根路径时展示login.html这个页面
		ctx.HTML(http.StatusOK, "login.html", gin.H{"client_key": clientKey})
	})

	engine.GET("/douyin_login", douyinLogin)

	// go run "D:/Program Files/Go/src/crypto/tls/generate_cert.go" --host="localhost"
	engine.RunTLS("localhost:5678", "config/keys/cert.pem", "config/keys/key.pem")
}
