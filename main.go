package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoExample holds internals for auth flow.
type CognitoExample struct {
	CognitoClient *cognito.CognitoIdentityProvider
	RegFlow       *regFlow
	UserPoolID    string
	AppClientID   string
}

type regFlow struct {
	Username string
}

func stream(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("./public%s.html", r.URL.Path)

	html, err := os.Open(path)
	defer html.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", path), http.StatusNotFound)
	}

	w.Header().Set("Content-type", "text/html")
	io.Copy(w, html)
}

func main() {
	conf := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}

	c := CognitoExample{
		CognitoClient: cognito.New(sess),
		RegFlow:       &regFlow{},
		UserPoolID:    os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:   os.Getenv("COGNITO_APP_CLIENT_ID"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			stream(w, r)
		case http.MethodPost:
			c.Register(w, r)
		}
	})

	http.HandleFunc("/otp", func(w http.ResponseWriter, r *http.Request) {

		if c.RegFlow.Username == "" {
			http.Redirect(w, r, "/register?error=You must register before sending OTP.", http.StatusFound)
		}

		switch r.Method {
		case http.MethodGet:
			stream(w, r)
		case http.MethodPost:
			c.OTP(w, r)
		}
	})

	http.HandleFunc("/username", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			stream(w, r)
		case http.MethodPost:
			c.Username(w, r)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			stream(w, r)
		case http.MethodPost:
			c.Login(w, r)
		}
	})

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
