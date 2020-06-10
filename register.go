package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Register handles sign in scenario.
func (c *CognitoExample) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	phoneNumber := r.Form.Get("phone_number")

	user := &cognito.SignUpInput{
		Username: aws.String(username),
		Password: aws.String(password),
		ClientId: aws.String(c.AppClientID),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
		},
	}

	_, err := c.CognitoClient.SignUp(user)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/register?error=%s", err.Error()), http.StatusSeeOther)
		return
	}

	c.RegFlow.Username = username

	http.Redirect(w, r, "/otp", http.StatusFound)

}
