package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

type ClaimsPage struct {
	AccessToken  string
	Claims       jwt.MapClaims
	RefreshToken string
}

var (
	clientID     string
	clientSecret string
	redirectURL  string
	issuerURL    string
	domain       string
	responseType = "code" // or "token" for implicit flow
	scope        = []string{"openid", "email", "phone"}
	signinURL    = "https://%s/login?client_id=%s&response_type=%s&scope=%s&%s"
	signupURL    = "https://%s/signup?client_id=%s&response_type=%s&scope=%s&%s"
	oauth2Config oauth2.Config
)

func init() {
	clientID = os.Getenv("COGNITO_CLIENT_ID")
	clientSecret = os.Getenv("COGNITO_CLIENT_SECRET")
	redirectURL = os.Getenv("COGNITO_REDIRECT_URL")
	issuerURL = os.Getenv("COGNITO_ISSUER_URL")
	domain = os.Getenv("COGNITO_DOMAIN")

	if clientID == "" || domain == "" || redirectURL == "" || issuerURL == "" {
		log.Fatal("Missing required environment variables: COGNITO_CLIENT_ID, COGNITO_DOMAIN, COGNITO_REDIRECT_URL, COGNITO_ISSUER_URL")
	}

	// Initialize OIDC provider
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}

	// Set up OAuth2 config
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scope,
	}
}

func main() {
	app := fiber.New()
	app.Get("/", handleLoginFiber)
	app.Get("/login", handleLoginFiber)
	app.Get("/signup", handleSignupFiber)
	app.Get("/callback", handleCallbackFiber)
	app.Get("/logout", handleLogoutFiber)
	app.Get("/cognito", handleCognitoFiber)

	fmt.Println("Fiber server is running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

func handleLoginFiber(c *fiber.Ctx) error {
	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	loginURL := fmt.Sprintf(signinURL, domain, clientID, responseType, strings.Join(scope, "+"), params.Encode())
	return c.Redirect(loginURL, fiber.StatusFound)
}

func handleSignupFiber(c *fiber.Ctx) error {
	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	signupURL := fmt.Sprintf(signupURL, domain, clientID, responseType, strings.Join(scope, "+"), params.Encode())
	return c.Redirect(signupURL, fiber.StatusFound)
}

func handleCallbackFiber(c *fiber.Ctx) error {
	ctx := context.Background()
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Code not found in query parameters")
	}

	rawToken, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token: " + err.Error())
	}

	tokenString := rawToken.AccessToken
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid claims")
	}

	pageData := ClaimsPage{
		AccessToken:  tokenString,
		Claims:       claims,
		RefreshToken: rawToken.RefreshToken,
	}

	tmpl, err := template.ParseFiles("../templates/claims.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error loading template: " + err.Error())
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, pageData); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template: " + err.Error())
	}

	return c.Type("html").SendString(buf.String())
}

func handleLogoutFiber(c *fiber.Ctx) error {
	return c.Redirect("/", fiber.StatusFound)
}

func handleCognitoFiber(c *fiber.Ctx) error {
	return c.SendString("Cognito callback handled")
}
