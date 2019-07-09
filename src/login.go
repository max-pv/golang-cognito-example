package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

// Login handles login scenario.
func (c *CognitoExample) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	refresh := r.Form.Get("refresh")
	refreshToken := r.Form.Get("refresh_token")

	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	if refresh != "" {
		flow = aws.String(flowRefreshToken)
		params = map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		}
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(c.AppClientID),
	}

	res, err := c.CognitoClient.InitiateAuth(authTry)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/login?error=%s", err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/login?authres=%s", res.AuthenticationResult), http.StatusFound)
}
