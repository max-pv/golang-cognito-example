package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Username handles username scenario.
func (c *CognitoExample) Username(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")

	_, err := c.CognitoClient.AdminGetUser(&cognito.AdminGetUserInput{
		UserPoolID: aws.String(c.UserPoolID),
		Username:   aws.String(username),
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok {
			if awsErr.Code() == cognito.ErrCodeUserNotFoundException {
				m := fmt.Sprintf("Username %s is free!", username)
				http.Redirect(w, r, fmt.Sprintf("/username?error=%s", m), http.StatusSeeOther)
				return
			}
		} else {
			http.Redirect(w, r, "Something went wrong", http.StatusSeeOther)
			return
		}
	}

	m := fmt.Sprintf("Username %s is taken.", username)
	http.Redirect(w, r, fmt.Sprintf("/username?error=%s", m), http.StatusSeeOther)
}
