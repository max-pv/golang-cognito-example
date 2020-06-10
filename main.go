package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

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
		UserPoolID:    os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:   os.Getenv("COGNITO_APP_CLIENT_ID"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect all requests to root to index.html
		if r.URL.Path == "/" {
			r.URL.Path = "/index"
		}

		handler := strings.Title(strings.ReplaceAll(r.URL.Path, "/", ""))
		// fmt.Printf("%#v\n", reflect.ValueOf(c).MethodByName("Login"))

		switch r.Method {
		case http.MethodGet:
			stream(w, r)
		case http.MethodPost:
			reflect.ValueOf(&c).
				MethodByName(handler).
				Call([]reflect.Value{
					reflect.ValueOf(w),
					reflect.ValueOf(r),
				})
		}
	})

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Starting Cognito Example on localhost%s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err.Error())
	}
}
