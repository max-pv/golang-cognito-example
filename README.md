# Golang AWS Cognito Register, Verify phone number, Login and Get User example

This example code demonstrates how to use AWS Cognito with AWS Go SDK in a form of simple web pages where you can:

1. Check if username is taken
2. Register
3. Verify user's phone
4. Login with username or refresh token

In order this solution to work, you need to have AWS credentials configured (file `.aws/configuration` exists) and User Pool created in AWS Console. You have to disable "Remember device" and enable "Sms second-factor" on authentication tab.

You will also need to create App Client in User Pool without "Generate Secret Key" checkbox. When the app client is created, in it's settings select "Enable username-password (non-SRP) flow for app-based authentication (USER_PASSWORD_AUTH)".

## Build

```go
go build -o ./build/cognito

AWS_PROFILE=XXX COGNITO_APP_CLIENT_ID=XXX COGNITO_USER_POOL_ID=XXX PORT=8080 ./build/cognito
```

Visit http://localhost:8080/register to see the registration page.