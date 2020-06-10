package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoExample holds internals for auth flow.
type CognitoExample struct {
	CognitoClient *cognito.CognitoIdentityProvider
	UserPoolID    string
	AppClientID   string
}

var routes = map[string]string{
	"/login":    "Login",
	"/otp":      "OTP",
	"/register": "Register",
	"/username": "Username",
}

// Stream responds with static HTML file
func Stream(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("./public%s.html", r.URL.Path)

	html, err := os.Open(path)
	defer html.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", path), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "text/html")
	io.Copy(w, html)
}

// Call Dynamically calls method of CognitoExample by name
func Call(c *CognitoExample, w http.ResponseWriter, r *http.Request) {
	handler, ok := routes[r.URL.Path]
	if !ok {
		http.Error(w, fmt.Sprintf("Handler %s not found", r.URL.Path), http.StatusNotFound)
		return
	}

	reflect.ValueOf(c).
		MethodByName(handler).
		Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
		})
}

func main() {
	conf := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	c := CognitoExample{
		CognitoClient: cognito.New(sess),
		UserPoolID:    os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:   os.Getenv("COGNITO_APP_CLIENT_ID"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect all requests to root to index.html
		if r.URL.Path == "/" {
			r.URL.Path = "/index"
		}

		switch r.Method {
		case http.MethodGet:
			// respond with static file
			Stream(w, r)
		case http.MethodPost:
			// dynamically call methods of CognitoExample
			// /login -> c.Login(w, r)
			// /register -> c.Register(w, r)
			Call(&c, w, r)
		}
	})

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Starting Cognito Example on localhost%s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
