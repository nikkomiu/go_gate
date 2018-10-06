package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

var domain string
var publicKey *rsa.PublicKey
var jwksData *jwksKey

// Auth0User to be retrieved from claims on successful JWT auth
type Auth0User struct {
	Email         string
	EmailVerified bool
	Picture       string
	Auth0ID       string
	UpdatedAt     string
}

type jwksKey struct {
	KeyID       string   `json:"kid"`
	Algorithm   string   `json:"alg"`
	Certificate []string `json:"x5c"`
}

type jwksResponse struct {
	Keys []*jwksKey `json:"keys"`
}

// Setup fetches the JWKS and builds up JWT auth
func Setup(authSettings *settings.AuthSettings) {
	client := &http.Client{}
	resp, err := client.Get(authSettings.JWKSURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data jwksResponse
	json.Unmarshal(body, &data)

	domain = authSettings.Domain
	jwksData = data.Keys[0]
	publicKey = convertKey(jwksData.Certificate[0])
}

// Validate to determine if user is authenticated
func Validate(authToken string, authRequred bool) (*Auth0User, error) {
	// Get the token from the Authorization header and remove Bearer
	if authToken != "" {
		authToken = strings.Split(authToken, " ")[1]
	}

	// Parse the token
	token, err := jwt.Parse(authToken, parseToken)
	if (err != nil || !token.Valid) && authRequred {
		return nil, errors.New("Authorization failed")
	}

	// Assign the Auth0User to the request settings
	user, err := parseUser(token)
	if err != nil {
		return nil, errors.New("Authorization failed")
	}

	// TODO: Add/Update User in DB (in a goroutine)

	return user, nil
}

func parseToken(token *jwt.Token) (interface{}, error) {
	// If the token signer or alg is wrong reject the token
	_, ok := token.Method.(*jwt.SigningMethodRSA)
	if !ok || token.Header["alg"] != jwksData.Algorithm || token.Header["kid"] != jwksData.KeyID {
		return nil, fmt.Errorf("Unexpected signing method")
	}

	// Return the successfully converted key
	return publicKey, nil
}

func parseUser(token *jwt.Token) (*Auth0User, error) {
	if r := recover(); r != nil {
		claims := token.Claims.(jwt.MapClaims)

		// Validate the issuer (if there is no issuer the token is invalid)
		if claims["iss"] != domain {
			return nil, errors.New("Invalid Issuer")
		}

		user := &Auth0User{
			Email:         claims["email"].(string),
			EmailVerified: claims["email_verified"].(bool),
			Picture:       claims["picture"].(string),
			Auth0ID:       strings.Split(claims["sub"].(string), "|")[1],
			UpdatedAt:     claims["updated_at"].(string),
		}

		return user, nil
	}

	return nil, nil
}

func convertKey(key string) *rsa.PublicKey {
	certPEM := "-----BEGIN CERTIFICATE-----\n" + key + "\n-----END CERTIFICATE-----"
	block, _ := pem.Decode([]byte(certPEM))
	cert, _ := x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	return rsaPublicKey
}
