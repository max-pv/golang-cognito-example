# Golang AWS Cognito Register, Verify phone number, Login and Get User example

## Just Show Me

If you are just curious how things work all together, you can find this example working at https://golang-cognito-example.herokuapp.com

## Instructions

This example code demonstrates how to use AWS Cognito with AWS Go SDK in a form of simple web pages where you can:

1. Check if username is taken
2. Register
3. Verify user's phone
4. Login with username or refresh token

In order this solution to work, you need to have AWS credentials configured (file `.aws/configuration` exists) and User Pool created in AWS Console. You have to disable "Remember device" and enable "Sms second-factor" on authentication tab.

When the app client is created, in it's settings select "Enable username-password (non-SRP) flow for app-based authentication (USER_PASSWORD_AUTH)".

It's possible to use go sdk with client secret. You can read a bit more about generating client secrets here:
https://dev.to/mcharytoniuk/using-aws-cognito-app-client-secret-hash-with-go-8ld

## Build

```go
go build -o ./build/cognito
```

## Run

Without client secret:

```go
AWS_PROFILE=XXX COGNITO_APP_CLIENT_ID=XXX COGNITO_USER_POOL_ID=XXX PORT=8080 ./build/cognito
```

With client secret:

```go
AWS_PROFILE=XXX COGNITO_APP_CLIENT_ID=XXX COGNITO_APP_CLIENT_SECRET=XXX  COGNITO_USER_POOL_ID=XXX PORT=8080 ./build/cognito
```

It's worth noting that in production environment you should not pass client secrets this way because with adequate permissions it's possible to read environmental variables of a running process. Also if you call a command that way, secret hash will be stored in your shell history. You should keep those issues in mind and mitigate them in your enviroment.

Visit http://localhost:8080/ to see the list of available pages.
