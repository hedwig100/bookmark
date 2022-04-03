package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/hedwig100/bookmark/backend/slog"
)

var pub *rsa.PublicKey
var pri *rsa.PrivateKey

func init() {
	// root indicates 'bookmark/backend'
	root := os.Getenv("BOOKMARK_ROOT")

	// prepare private key
	buf, err := ioutil.ReadFile(fmt.Sprintf("%sprivate.pem", root))
	if err != nil {
		slog.Fatalf("cannot open private key : %v", err)
	}

	block, _ := pem.Decode(buf)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		slog.Fatalf("not private key : %v", err)
	}

	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		slog.Fatalf("cannot parse private key : %v", err)
	}

	slog.Info("Preparing private-key successful!")

	// prepare public key
	buf, err = ioutil.ReadFile(fmt.Sprintf("%spublic.pem", root))
	if err != nil {
		slog.Fatalf("cannot open public key : %v", err)
	}

	block, _ = pem.Decode(buf)
	if block == nil || block.Type != "PUBLIC KEY" {
		slog.Fatalf("not public key : %v", err)
	}

	pub_, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		slog.Fatalf("cannot parse public key : %v", err)
	}

	var ok bool
	pub, ok = pub_.(*rsa.PublicKey)
	if !ok {
		slog.Fatalf("not rsa public key : %v", err)
	}

	slog.Info("Preparing public key successful!")
}

// GenJWT generates a JWT and returns it to the client.
func GenJWT(w http.ResponseWriter, user_id string, username string) {
	// generate header and claims
	// "iss" (issuer): the principal that issued the JWT.
	// "sub" (subject): the principal that is the subject of the JWT.
	// "exp" (expiration time): the expiration time on or after which the JWT MUST NOT be accepted for processing
	// "iat" (issued at): the time at which the JWT was issued.
	// about iat https://github.com/dgrijalva/jwt-go/issues/383
	// TODO: fix iat problem
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":      "https://localhost",
		"sub":      user_id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"username": username,
	})
	tokenStr, err := token.SignedString(pri)
	if err != nil {
		slog.Errf("internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// response in cookie
	c := http.Cookie{
		Name:     "bookmark_auth",
		Value:    tokenStr,
		MaxAge:   60 * 60 * 24,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
}

// Auth verifies credentials using Authorization in the request header and a private key.
// If not valid, it returns a 403 error.
func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// we expect "Authorization: Bearer <token>"
		c, err := r.Cookie("bookmark_auth")
		if err != nil {
			slog.Info("Authorization isn't set")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenStr := c.Value

		// parse jwt
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			} else {
				return pub, nil
			}
		})
		slog.Debug(token)
		slog.Debug(token.Claims)

		if err != nil {
			slog.Infof("Unauthorized: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err = token.Claims.Valid(); err != nil || !token.Valid {
			slog.Infof("Unauthorized: %v or token isn't valid.", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// check if correct user
		params := httptreemux.ContextParams(r.Context())
		username, ok := params["username"]
		if !ok {
			slog.Err("Route must have 'username' in the path.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Err("Claim type must be jwt.MapClaims.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claimUsername, ok := claim["username"]
		if !ok {
			slog.Infof("Unauthorized: there is no 'username' key in claim.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != claimUsername {
			slog.Infof("Unauthorized: unexpected username")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// handler
		handler(w, r)
	}
}
