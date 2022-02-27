package services

//go:generate mockgen -source=$GOFILE -destination=../mocks/services/mock_$GOFILE -package=mockservices

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IAuth interface {
	ValidateSession(tokenString string) (*jwt.Token, error)
	Authenticate(username, password string) (bool, error)
	CreateAuthToken(username string) (string, error)
}

type Auth struct {
	JWTIssuer     string
	SigningMethod jwt.SigningMethod
	PrivateKey    *rsa.PrivateKey
	UserAccounts  repositories.IUserAccounts
}

func (a *Auth) ValidateSession(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		switch t.Method.Alg() {
		case a.SigningMethod.Alg():
			return a.PrivateKey.Public(), nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *Auth) Authenticate(username, password string) (bool, error) {
	userExists, err := a.UserAccounts.ExistsByUsername(username)
	if err != nil {
		return false, fmt.Errorf("error checking if user exists: %w", err)
	}
	if !userExists {
		return false, constants.ErrUserDoesNotExist
	}

	userAccount, err := a.UserAccounts.GetByUsername(username)
	if err != nil {
		return false, fmt.Errorf("error getting user account: %w", err)
	}

	isValid, err := passwordIsValid(password, userAccount.PasswordHash)
	if err != nil {
		return false, fmt.Errorf("error checking password validity: %w", err)
	}

	return isValid, nil
}

func (a *Auth) CreateAuthToken(username string) (string, error) {
	issueTime := time.Now()
	token := jwt.NewWithClaims(a.SigningMethod, jwt.StandardClaims{
		Issuer:    a.JWTIssuer,
		Audience:  a.JWTIssuer,
		Subject:   username,
		IssuedAt:  issueTime.Unix(),
		ExpiresAt: issueTime.Add(60 * time.Minute).Unix(),
	})

	signedTokenString, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return signedTokenString, nil
}

func passwordIsValid(password string, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
