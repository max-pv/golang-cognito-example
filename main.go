package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/br4in3x/golang-cognito-example/app"
)

// Stream responds with static HTML file
func Stream(w http.ResponseWriter, r *http.Request) {
	// Redirect all requests to root to index.html
	if r.URL.Path == "/" {
		r.URL.Path = "/index"
	}

	path := fmt.Sprintf("./static%s.html", r.URL.Path)

	html, err := os.Open(path)
	defer html.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", path), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "text/html")
	io.Copy(w, html)
}

// Call routes POST requests
func Call(a *app.App, w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {
	case "/login":
		a.Login(w, r)
	case "/otp":
		a.OTP(w, r)
	case "/register":
		a.Register(w, r)
	case "/username":
		a.Username(w, r)
	default:
		http.Error(w, fmt.Sprintf("Handler for POST %s not found", r.URL.Path), http.StatusNotFound)
	}

	return
}

func main() {
	conf := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	example := app.App{
		CognitoClient:   cognito.New(sess),
		UserPoolID:      os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:     os.Getenv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: os.Getenv("COGNITO_APP_CLIENT_SECRET"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Respond with static file
			Stream(w, r)
		case http.MethodPost:
			// Handle html form submission
			Call(&example, w, r)
		}
	})

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Starting Cognito Example on localhost%s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
