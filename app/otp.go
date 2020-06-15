package app

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// OTP handles phone verification step.
func (a *App) OTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	otp := r.Form.Get("otp")
	username := r.Form.Get("username")

	user := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(otp),
		Username:         aws.String(username),
		ClientId:         aws.String(a.AppClientID),
	}

	_, err := a.CognitoClient.ConfirmSignUp(user)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/otp?message=%s", err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/username", http.StatusFound)
}
