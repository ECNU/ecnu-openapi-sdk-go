package sdk

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	userInfoURL string
	config      *oauth2.Config
	client      *http.Client
	authLock    = new(sync.RWMutex)
	c           *cache.Cache
)

type UserInfoResponse struct {
	ErrCode   int    `json:"errCode"`
	ErrMsg    string `json:"errMsg"`
	RequestId string `json:"requestId"`
	Data      struct {
		UserId     string `json:"userId"`
		Name       string `json:"name"`
		VpnEnabled int    `json:"vpnEnabled"`
	} `json:"data"`
}

func InitOAuth2AuthorizationCode(cf OAuth2Config) {
	scopes := []string{DefaultScope}
	authURL := DefaultAuthURL
	tokenURL := DefaultTokenURL
	userInfoURL = DefaultUserInfoURL
	if len(cf.Scopes) > 0 {
		scopes = cf.Scopes
	}
	if cf.UserInfoURL != "" {
		userInfoURL = cf.UserInfoURL
	}
	if cf.Endpoint.AuthURL != "" {
		authURL = cf.Endpoint.AuthURL
	}
	if cf.Endpoint.TokenURL != "" {
		tokenURL = cf.Endpoint.TokenURL
	}

	config = &oauth2.Config{
		ClientID:     cf.ClientId,
		ClientSecret: cf.ClientSecret,
		RedirectURL:  cf.RedirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	var expiration time.Duration
	if cf.Cache.Expiration == 0 {
		expiration = DefaultCacheExpiration
	} else {
		expiration = cf.Cache.Expiration
	}
	var cleanup time.Duration
	if cf.Cache.Cleanup == 0 {
		cleanup = DefaultCacheCleanup
	} else {
		cleanup = cf.Cache.Cleanup
	}
	c = cache.New(expiration, cleanup)

}

func GetAuthorizationEndpoint(state string) string {
	c.Set(state, "", cache.DefaultExpiration)
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func GetInfo(code string, state string) (UserInfoResponse, error) {
	token, err := GetToken(code, state)
	if err != nil {
		return UserInfoResponse{}, err
	}
	client := GetClient(token)
	return GetUserInfo(client)
}

func GetToken(code string, state string) (*oauth2.Token, error) {
	_, found := c.Get(state)
	if !found {
		return nil, fmt.Errorf("state有误")
	}
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("获取token失败")
		return nil, err
	}
	return token, nil
}

func GetUserInfo(client *http.Client) (UserInfoResponse, error) {
	resp, err := client.Get(userInfoURL)
	if err != nil {
		fmt.Printf("Get请求失败: %v", err)
		return UserInfoResponse{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭响应体失败")
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体失败: %v\n", err)
		return UserInfoResponse{}, err
	}

	var userInfo UserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return UserInfoResponse{}, err
	}

	return userInfo, nil
}

func GetClient(token *oauth2.Token) *http.Client {
	authLock.RLock()
	client = config.Client(oauth2.NoContext, token)
	defer authLock.RUnlock()
	return client
}

func GenerateState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("生成随机state失败: %v", err)
	}
	return hex.EncodeToString(b)
}
