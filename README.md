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
1. Clone the repository:
   ```sh
   git clone https://github.com/ksquare-mchavez/ksquare-biolinq-cognito.git
   cd ksquare-biolinq-cognito
   ```
2. Install dependencies:
   ```sh
   go get ./...
   ```
3. Update Cognito credentials in `main.go`:
   - `clientID`, `clientSecret`, `redirectURL`, `issuerURL`
4. Run the server:
   ```sh
   go run main.go
   ```
5. Open [http://localhost:8080](http://localhost:8080) in your browser.

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
