package main

import (
	"fmt"
	"github.com/ecnu/ecnu-openapi-sdk-go/sdk"
	"net/http"
)

// 定义OAuth2配置
var oauth2Config = sdk.OAuth2Config{
	ClientId:     "client_id",
	ClientSecret: "client_secret",
	RedirectURL:  "http://localhost:8080/user",
	Scopes:       []string{"ECNU-Basic"},
	Endpoint: sdk.EndpointConf{
		AuthURL:  "https://api.ecnu.edu.cn/oauth2/authorize",
		TokenURL: "https://api.ecnu.edu.cn/oauth2/token",
	},
}

func main() {

	http.HandleFunc("/login", login)

	http.HandleFunc("/user", getUserInfo)

	sdk.InitOAuth2AuthorizationCode(oauth2Config)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe Error: ", err)
		return
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	state := sdk.GenerateState()

	url := sdk.GetAuthorizationEndpoint(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" {
		fmt.Fprintf(w, "Code not found")
		return
	}
	userInfo, err := sdk.GetInfo(code, state)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>User Info</h1>")
	fmt.Fprintf(w, "<p>User ID: %s</p>", userInfo.Data.UserId)
	fmt.Fprintf(w, "<p>Name: %s</p>", userInfo.Data.Name)
	vpnStatus := "Disabled"
	if userInfo.Data.VpnEnabled == 1 {
		vpnStatus = "Enabled"
	}
	fmt.Fprintf(w, "<p>VPN Enabled: %s</p>", vpnStatus)
	fmt.Fprintf(w, "</body></html>")
}
