package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	MAXPageSIZE = 10000
)

/*
APIResult 数据响应结构
https://developer.ecnu.edu.cn/doc/#/architecture/design?id=%e6%95%b0%e6%8d%ae%e5%93%8d%e5%ba%94%e7%bb%93%e6%9e%84
*/
type APIResult struct {
	ErrCode   int64           `json:"errCode"`
	ErrMsg    string          `json:"errMsg"`
	RequestId string          `json:"requestId"`
	Data      json.RawMessage `json:"data"`
}

func parseApiResult(result *http.Response, debug bool) (APIResult, error) {
	var data APIResult
	if debug {
		fmt.Println(result.Request.URL.String())
		fmt.Println(result.StatusCode)
		fmt.Println(result.Header)
	}
	if result.StatusCode != 200 {
		//错误码：A401OT access_token 参数错误。清空 access_token 再来一次
		if result.Header.Get("X-Ca-Error-Code") == "A401OT" {
			return data, errors.New("A401OT")
		}
		return data, fmt.Errorf("invoke api get fail, X-Ca-Error-Code: %s, X-Ca-Error-Message: %s,X-Ca-Request-Id: %s",
			result.Header.Get("X-Ca-Error-Code"),
			result.Header.Get("X-Ca-Error-Message"),
			result.Header.Get("X-Ca-Request-Id"))
	}

	if result.Body == nil {
		return data, fmt.Errorf("get api response body is nil")
	}

	defer result.Body.Close()
	res, err := io.ReadAll(result.Body)
	if debug {
		fmt.Println(string(res))
	}
	if err != nil {
		return data, fmt.Errorf("read getapi response body fail: %v", err)
	}

	if err = json.Unmarshal(res, &data); err != nil {
		return data, fmt.Errorf("parse getapi response body fail: %v", err)
	}
	return data, nil
}

// HttpGet 通用GET请求
func (c *OAuth2Client) HttpGet(url string) (json.RawMessage, error) {
	var apiResult APIResult
	var err error

	result, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("invoke api get fail: %v", err)
	}
	apiResult, err = parseApiResult(result, c.Debug)
	if err != nil {
		if err.Error() == "A401OT" && c.RetryCount <= 3 {
			//错误码：A401OT access_token 参数错误。清空再来一次
			c.retryAdd()
			c.Client = c.conf.Client(context.Background())
			return c.HttpGet(url)
		}
		return nil, err
	}
	if c.RetryCount > 0 {
		//token 是正确的，如果之前有计数，这里要清零了
		c.retryRest()
	}
	if apiResult.ErrCode != 0 {
		return nil, errors.New(apiResult.ErrMsg)
	}

	return apiResult.Data, nil
}

// 通用 http 请求
// todo
func (c *OAuth2Client) HttpRequest(url, method string, header map[string]string, body io.Reader) (json.RawMessage, error) {
	return nil, nil
}

func CallAPI(url, method string, header map[string]string, body io.Reader) (json.RawMessage, error) {
	c := GetOpenAPIClient()
	switch method {
	case "GET":
		return c.HttpGet(url)
	}
	return nil, errors.New("not support method")
}
