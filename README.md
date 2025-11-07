# KSquare Biolinq Cognito

A Go web application using Fiber and AWS Cognito for OIDC authentication.

## Features
- Login with AWS Cognito (OIDC)
- JWT claims display
- Modern UI with HTML templates

## Requirements
- Go 1.18+
- AWS Cognito User Pool


## Example curl Commands

### Test Login Redirect
```sh
curl -v http://localhost:8080/login -L
```
This will follow the redirect to Cognito's login page.

### Test Callback (after login)
```sh
curl -v "http://localhost:8080/callback?code=YOUR_AUTH_CODE"
```
Replace `YOUR_AUTH_CODE` with the code received from Cognito after login.

### Test Home Page
```sh
curl -v http://localhost:8080/
```

## Project Structure
```
├── main.go
├── go.mod
├── templates
│   ├── home.html
│   └── claims.html
└── README.md
```

## License
MIT
