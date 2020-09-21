package app

import (
	// Those imports are required to compute the secret hash.
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"

	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

// Secret hash is not a client secret itself, but a base64 encoded hmac-sha256
// hash.
// The actual AWS documentation on how to compute this hash is here:
// https://docs.aws.amazon.com/cognito/latest/developerguide/signing-up-users-in-your-app.html#cognito-user-pools-computing-secret-hash
func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Login handles login scenario.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
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

	// Compute secret hash based on client secret.
	if a.AppClientSecret != "" {
		secretHash := computeSecretHash(a.AppClientSecret, username, a.AppClientID)

		params["SECRET_HASH"] = aws.String(secretHash)
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
		ClientId:       aws.String(a.AppClientID),
	}

	res, err := a.CognitoClient.InitiateAuth(authTry)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/login?message=%s", err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/login?authres=%s", res.AuthenticationResult), http.StatusFound)
}
