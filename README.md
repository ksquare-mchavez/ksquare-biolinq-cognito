# KSquare Biolinq Cognito

A Go web application using Fiber and AWS Cognito for OIDC authentication.

## Features
- Login with AWS Cognito (OIDC)
- JWT claims display
- Modern UI with HTML templates

## Requirements
- Go 1.18+
- AWS Cognito User Pool

## Setup

1. **Clone the repository**
   ```sh
   git clone https://github.com/biolinq-inc/ksquare-biolinq-cognito.git
   cd ksquare-biolinq-cognito
   ```

2. **Install dependencies**
   ```sh
   go mod tidy
   ```

3. **Configure AWS Cognito credentials**
   Set Cognito credentials as environment variables:
	 ```
   export COGNITO_CLIENT_ID=your_client_id
   export COGNITO_CLIENT_SECRET=your_client_secret
   export COGNITO_REDIRECT_URL=http://localhost:8080/callback
   export COGNITO_ISSUER_URL=https://cognito-idp.us-east-2.amazonaws.com/your_pool_id
	 ```

4. **Run the server**
   ```sh
   go run cmd/main.go
   ```

   The application will start on `http://localhost:8080`.


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
