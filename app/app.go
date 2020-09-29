package app

import (
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// App holds internals for auth flow.
type App struct {
	CognitoClient   *cognito.CognitoIdentityProvider
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}
