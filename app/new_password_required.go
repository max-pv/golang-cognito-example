package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// NewPasswordRequired handles set new password.
func (a *App) NewPasswordRequired(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	session := r.Form.Get("session")

	params := map[string]*string{
		"USERNAME":     aws.String(username),
		"NEW_PASSWORD": aws.String(password),
	}
	// Compute secret hash based on client secret.
	if a.AppClientSecret != "" {
		secretHash := computeSecretHash(a.AppClientSecret, username, a.AppClientID)

		params["SECRET_HASH"] = aws.String(secretHash)
	}
	chall := &cognito.RespondToAuthChallengeInput{
		ChallengeName:      aws.String("NEW_PASSWORD_REQUIRED"),
		ClientId:           aws.String(a.AppClientID),
		ChallengeResponses: params,
		Session:            aws.String(session),
	}
	res, err := a.CognitoClient.RespondToAuthChallenge(chall)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/login?message=%s", err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/login?authres=%s", res.AuthenticationResult), http.StatusFound)
}
