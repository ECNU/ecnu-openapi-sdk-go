package sdk

import (
	"context"
	"net/http"
	"sync"
	"time"

	cc "golang.org/x/oauth2/clientcredentials"
)

const (
	DefaultScope   = "ECNU-Basic"
	DefaultBaseURL = "https://api.ecnu.edu.cn"
	DefaultTimeout = 10
)

var (
	openAPIClient *OAuth2Client
	lock          = new(sync.RWMutex)
)

type OAuth2Client struct {
	conf       *cc.Config
	Client     *http.Client
	BaseUrl    string
	RetryCount int
	Debug      bool
}

type OAuth2Config struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	BaseUrl      string   `json:"base_url"`
	Scopes       []string `json:"scopes"`
	Timeout      int64    `json:"timeout"`
	Debug        bool     `json:"debug"`
}

// Init 初始化 OAuth2 应用
func InitOAuth2ClientCredentials(cf OAuth2Config) {
	baseUrl := DefaultBaseURL
	scopes := []string{DefaultScope}
	var timeout int64 = DefaultTimeout

	if cf.BaseUrl != "" {
		baseUrl = cf.BaseUrl
	}
	if len(cf.Scopes) > 0 {
		scopes = cf.Scopes
	}
	if cf.Timeout > 0 {
		timeout = cf.Timeout
	}
	conf := &cc.Config{
		ClientID:     cf.ClientId,
		ClientSecret: cf.ClientSecret,
		Scopes:       scopes,
		TokenURL:     baseUrl + "/oauth2/token",
	}
	client := conf.Client(context.Background())
	client.Timeout = time.Second * time.Duration(timeout)

	openAPIClient = &OAuth2Client{Client: client, BaseUrl: baseUrl, Debug: cf.Debug}
}

// GetOpenAPIClient 获取接口的Client信息
func GetOpenAPIClient() *OAuth2Client {
	lock.RLock()
	defer lock.RUnlock()
	return openAPIClient
}

func (c *OAuth2Client) retryAdd() {
	lock.RLock()
	defer lock.RUnlock()
	c.RetryCount = c.RetryCount + 1
}

func (c *OAuth2Client) retryRest() {
	lock.RLock()
	defer lock.RUnlock()
	c.RetryCount = 0
}
