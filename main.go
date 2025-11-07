package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/coreos/go-oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

type ClaimsPage struct {
	AccessToken string
	Claims      jwt.MapClaims
}

var (
	clientID     = "33vu9td0v0l3sll8sva5vu8j72"
	clientSecret = "1pvnab10o7vo8avuk65sqmhjggk7729sgvtnr3od98j00d3rrd2i"
	redirectURL  = "http://localhost:8080/callback"
	issuerURL    = "https://cognito-idp.us-east-2.amazonaws.com/us-east-2_MuLHvXh19"
	oauth2Config oauth2.Config
)

func init() {
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
		Scopes:       []string{"openid", "email", "phone"},
	}
}

func main() {
	app := fiber.New()
	app.Get("/", handleHomeFiber)
	app.Get("/login", handleLoginFiber)
	app.Get("/callback", handleCallbackFiber)
	app.Get("/logout", handleLogoutFiber)
	app.Get("/cognito", handleCognitoFiber)

	fmt.Println("Fiber server is running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

// Fiber route handlers
func handleHomeFiber(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error loading template: " + err.Error())
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error rendering template: " + err.Error())
	}

	return c.Type("html").SendString(buf.String())
}

func handleLoginFiber(c *fiber.Ctx) error {
	state := "state" // Replace with a secure random string in production
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusFound)
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
		AccessToken: tokenString,
		Claims:      claims,
	}

	tmpl, err := template.ParseFiles("templates/claims.html")
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
