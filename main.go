package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/br4in3x/golang-cognito-example/app"
)

func getStaticFilePath(path string) (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errors.New("can not get current executable path")
	}

	// Redirect all requests to root to index.html
	if path == "/" {
		path = "/index"
	}

	absFilePath := fmt.Sprintf("%s/static%s.html", filepath.Dir(ex), path)
	return absFilePath, nil
}

// Stream responds with static HTML file
func Stream(w http.ResponseWriter, r *http.Request) {
	htmlFile, err := getStaticFilePath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html, err := os.Open(htmlFile)
	defer html.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", htmlFile), http.StatusNotFound)
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
