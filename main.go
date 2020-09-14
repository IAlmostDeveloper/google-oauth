package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	googleOauthConfig = &.Config{
		RedirectURL: "http://localhost:8080/callback",
		ClientID: "574484498021-l0tqmv2jfa9t9akmsruutva6sd0824a6.apps.googleusercontent.com",
		ClientSecret:  "fkK3H-ovtu7pcTdOxzFojMmb",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/plus.me"},
		Endpoint: google.Endpoint,
	}
	randomState = "random"
)

func main(){
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleCallback(writer http.ResponseWriter, request *http.Request) {
	if request.FormValue("state") != randomState{
		fmt.Println("state is not valid")
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	token ,err := googleOauthConfig.Exchange(.NoContext, request.FormValue("code"))
	if err != nil{
		fmt.Printf("could not get token : %s \n", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil{
		fmt.Printf("could not create get request : %s \n", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Printf("could not parse response : %s \n", err.Error())
		http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(writer, "Response: %s", content)
}

func handleLogin(writer http.ResponseWriter, request *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(writer, request, url, http.StatusTemporaryRedirect)
}

func handleHome(writer http.ResponseWriter, request *http.Request) {
	var html = `<html><body><a href="/login">Google auth</a></body></html>`
	fmt.Fprint(writer,  html)
}