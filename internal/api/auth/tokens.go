package auth

import (
	"crypto/sha512"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/startdusk/finance-app-backend/internal/model"
)

// verify jwt key
var jwtKey = []byte("my_secret_key") // TODO: change key
var tokenDurationHours = 30 * 24     // 30 days

type Tokens interface {
	IssueToken(principal model.Principal) (string, error)
	Verify(token string) (*model.Principal, error)
}

type tokens struct {
	key             []byte
	duration        time.Duration
	beforeTolerance time.Duration
	signingMethod   jwt.SigningMethod
}

type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

func NewTokens() Tokens {
	hasher := sha512.New()

	if _, err := hasher.Write([]byte(jwtKey)); err != nil {
		panic(err)
	}

	tokenDuration := time.Duration(tokenDurationHours) * time.Hour
	return &tokens{
		key:             hasher.Sum(nil),
		duration:        tokenDuration,
		beforeTolerance: -2 * time.Minute,
		signingMethod:   jwt.SigningMethodHS512,
	}
}

func (t *tokens) IssueToken(principal model.Principal) (string, error) {
	if principal.UserID == model.NilUserID {
		return "", errors.New("invalid principal")
	}
	now := time.Now()
	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			NotBefore: now.Add(t.beforeTolerance).Unix(),
			ExpiresAt: now.Add(t.duration).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(t.signingMethod, claims)
	// Create the JWT string
	return token.SignedString(t.key)
}

func (t *tokens) Verify(token string) (*model.Principal, error) {
	// TODO
	return nil, nil
}
