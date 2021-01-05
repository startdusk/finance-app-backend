package auth

import (
	// "crypto/sha512"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/startdusk/finance-app-backend/internal/model"
)

// verify jwt key
var jwtKey = []byte("my_secret_key") // TODO: change key

var accessTokenDuration = time.Duration(30) * time.Minute // 30 min
// var accessTokenDuration = time.Duration(60) * time.Second // TEST: we will get this token, try to get user, first time we should get one. we will wait for 1 miniute and try again. it should fail.
var refreshTokenDuration = time.Duration(30*24) * time.Hour // 30 days
// var refreshTokenDuration = time.Duration(120) * time.Second // TEST: after access token fails we send this token and it should return new one. we will wait 1 miniute and we should get expried error

type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

// Tokens is wrapper for access and refresh tokens
type Tokens struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"expiresAt,omitempty"` // we return's only access token's expires at time
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"-"` // we will store this time in database with refresh token
}

// IssueToken generate Access and Refresh tokens
// Currently we will generate only access token
// Refresh token I will generate the same way like access token
// I changed my mind about refresh token format. It will be randomly generated string (32/64 length)
// I will do it with JWT format. because I want to have user ID inside token so we will know who user is. maybe we will want same additional information
func IssueToken(principal model.Principal) (*Tokens, error) {
	if principal.UserID == model.NilUserID {
		return nil, errors.New("invalid principal")
	}
	// Generate Access token
	accessToken, accessTokenExpiresAt, err := generateToken(principal, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	// Generate Refresh token
	refreshToken, refreshTokenExpiresAt, err := generateToken(principal, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}

	return &tokens, nil
}

func generateToken(principal model.Principal, duration time.Duration) (string, int64, error) {
	now := time.Now()
	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}

func VerifyToken(token string) (model.Principal, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return model.NilPrincipal, err
		}

		return model.NilPrincipal, err
	}

	principal := model.Principal{
		UserID: claims.UserID,
	}

	// we want to return principal even token invalid becase we need to get userID
	if !tkn.Valid {
		return model.NilPrincipal, err
	}

	return principal, nil
}
