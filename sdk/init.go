package sdk

import (
	"context"
	"net/http"
	"sync"
	"time"

	cc "golang.org/x/oauth2/clientcredentials"
)

const (
	DefaultScope       = "ECNU-Basic"
	DefaultBaseURL     = "https://api.ecnu.edu.cn"
	DefaultTimeout     = 10
	DefaultUserInfoURL = "https://api.ecnu.edu.cn/oauth2/userinfo"
	DefaultAuthURL     = "https://api.ecnu.edu.cn/oauth2/authorize"
	DefaultTokenURL    = "https://api.ecnu.edu.cn/oauth2/token"

	DefaultCacheExpiration = 5 * time.Minute
	DefaultCacheCleanup    = 10 * time.Minute
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

type EndpointConf struct {
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}
type CacheConfig struct {
	Expiration time.Duration
	Cleanup    time.Duration
}

type OAuth2Config struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	Debug        bool     `json:"debug"`

	BaseUrl string `json:"base_url"`
	Timeout int64  `json:"timeout"`

	RedirectURL string       `json:"redirect_url"`
	UserInfoURL string       `json:"user_info_url"`
	Endpoint    EndpointConf `json:"endpoint"`

	Cache CacheConfig `json:"cache"`
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
