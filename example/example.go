package main

import (
	"fmt"

	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
)

func main() {
	/*
	   type OAuth2Config struct {
	   	ClientId     string   `json:"client_id"`
	   	ClientSecret string   `json:"client_secret"`
	   	BaseUrl      string   `json:"base_url"`     // 默认 https://api.ecnu.edu.cn
	   	Scopes       []string `json:"scopes"`       //默认 ["ECNU-Basic"]
	   	Timeout      int64    `json:"timeout"`      //默认10秒
	   	Debug        bool     `json:"debug"`        //默认 false, 如果开启 debug，会打印出请求和响应的详细信息，对于数据同步类接口而言可能会非常大
	   }
	*/
	cf := sdk.OAuth2Config{
		ClientId:     "client_id",
		ClientSecret: "client_secret",
	}
	sdk.InitOAuth2ClientCredentials(cf)

	// 直接调用接口

	res, err := sdk.CallAPI("https://api.ecnu.edu.cn/api/v1/sync/fakewithts?ts=0&pageSize=1&pageNum=1", "GET", nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))

	exampleSyncToCSV()
	exampleSyncToModel()
	exampleSyncToDB()
	//exampleSyncPerformance()
}
