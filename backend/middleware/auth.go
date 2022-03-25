package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var pub *rsa.PublicKey
var pri *rsa.PrivateKey

func init() {

	// prepare private key
	buf, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Fatalf("cannot open private key : %v", err)
		os.Exit(1)
	}

	block, _ := pem.Decode(buf)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.Fatalf("not private key : %v", err)
		os.Exit(1)
	}

	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("cannot parse private key : %v", err)
		os.Exit(1)
	}

	log.Println("Preparing private-key successful!")

	// prepare public key
	buf, err = ioutil.ReadFile("public.pem")
	if err != nil {
		log.Fatalf("cannot open public key : %v", err)
		os.Exit(1)
	}

	block, _ = pem.Decode(buf)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatalf("not public key : %v", err)
		os.Exit(1)
	}

	pub_, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("cannot parse public key : %v", err)
		os.Exit(1)
	}

	var ok bool
	pub, ok = pub_.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("not rsa public key : %v", err)
		os.Exit(1)
	}

	log.Println("Preparing public key successful!")
}

// GenJWT generates a JWT and returns it to the client. Redirect here after user authentication.
// Method:
// 		GET
// URI:
//		/auth
func GenJWT(w http.ResponseWriter, r *http.Request) {
	// generate header and claims
	// "iss" (issuer): the principal that issued the JWT.
	// "sub" (subject): the principal that is the subject of the JWT.
	// "exp" (expiration time): the expiration time on or after which the JWT MUST NOT be accepted for processing
	// "iat" (issued at): the time at which the JWT was issued.
	// TODO: use valid "sub" depending on the user
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://localhost",
		"sub": 1,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now(),
	})
	tokenStr, err := token.SignedString(pri)
	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] internal error: %v", err))
	}

	// response
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenStr))
	w.WriteHeader(http.StatusOK)
}

// Auth verifies credentials using Authorization in the request header and a private key.
// If not valid, it returns a 403 error.
func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// we expect "Authorization: Bearer <token>"
		auth, ok := r.Header["Authorization"]
		if !ok || !strings.HasPrefix(auth[0], "Bearer ") {
			log.Println("Authorization isn't set")
			w.WriteHeader(http.StatusUnauthorized)
		}
		tokenStr := strings.TrimPrefix(auth[0], "Bearer ")

		// parse jwt
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			} else {
				return pub, nil
			}
		})
		log.Println(token)
		log.Println(token.Claims)

		if err != nil {
			log.Println("Unauthorized: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		if err = token.Claims.Valid(); err != nil || !token.Valid {
			log.Println("Unauthorized: %v or token isn't valid.", err)
			w.WriteHeader(http.StatusUnauthorized)
		}

		// handler
		handler(w, r)
	}
}
