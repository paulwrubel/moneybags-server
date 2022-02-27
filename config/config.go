package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulwrubel/moneybags-server/constants"
	log "github.com/sirupsen/logrus"
)

type AppInfo struct {
	DB       *pgxpool.Pool
	AuthInfo *AuthInfo
}

type dBInfo struct {
	Host     string
	Username string
	Password string
}

type AuthInfo struct {
	JWTIssuer     string
	SigningMethod jwt.SigningMethod
	PrivateKey    *rsa.PrivateKey
}

func InitializeApp() (*AppInfo, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("error initializing db connection: %w", err)
	}
	authInfo, err := getAuthInfo()
	if err != nil {
		return nil, fmt.Errorf("error initializing auth info: %w", err)
	}

	return &AppInfo{
		DB:       db,
		AuthInfo: authInfo,
	}, nil
}

func getDB() (*pgxpool.Pool, error) {
	log.Info("getting DB info")
	pgHost, isSet := os.LookupEnv(constants.PostgresHostnameEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresHostnameEnvironmentKey)
	}
	pgUser, isSet := os.LookupEnv(constants.PostgresUsernameEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresUsernameEnvironmentKey)
	}
	pgPass, isSet := os.LookupEnv(constants.PostgresPasswordEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresPasswordEnvironmentKey)
	}

	log.Debug("initializing database")

	// initialize configuration
	connectionString := fmt.Sprintf("host=%s user=%s password=%s", pgHost, pgUser, pgPass)
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// initialize connection pool
	connectionAttempts := 0
	var db *pgxpool.Pool
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		db, err = pgxpool.ConnectConfig(ctx, poolConfig)
		if err == nil {
			break
		}
		connectionAttempts++
		if connectionAttempts >= 10 {
			return nil, err
		}
		cancel()
		// retry db
		log.WithError(err).Error("database connection attempt failed, waiting 5s then retrying")
		time.Sleep(time.Second * 5)
	}

	log.Debug("database initialized")
	return db, nil
}

func getAuthInfo() (*AuthInfo, error) {
	log.Info("getting auth info")
	// parse locations from environment variables
	jwtIssuer, isSet := os.LookupEnv(constants.JWTIssuerEnvironmentKey)
	if !isSet {
		jwtIssuer = constants.DefaultJWTIssuer
	}
	jwtRSAPrivateKeyFile, isSet := os.LookupEnv(constants.JWTRSAPrivateKeyFileEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.JWTRSAPrivateKeyFileEnvironmentKey)
	}
	jwtSigningAlg, isSet := os.LookupEnv(constants.JWTSigningAlgorithmEnvironmentKey)
	if !isSet {
		jwtSigningAlg = constants.DefaultJWTSigningAlgorithm
	}

	privateKeyBytes, err := ioutil.ReadFile(jwtRSAPrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %s", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %s", err)
	}

	return &AuthInfo{
		JWTIssuer:     jwtIssuer,
		SigningMethod: jwt.GetSigningMethod(jwtSigningAlg),
		PrivateKey:    privateKey,
	}, nil
}
