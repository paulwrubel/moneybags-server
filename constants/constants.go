package constants

import "errors"

const (
	PostgresHostnameEnvironmentKey = "MONEYBAGS_PG_HOST"
	PostgresUsernameEnvironmentKey = "MONEYBAGS_PG_USER"
	PostgresPasswordEnvironmentKey = "MONEYBAGS_PG_PASS"
)

const (
	// sensible default for jwt info
	DefaultJWTIssuer           = "moneybags"
	DefaultJWTSigningAlgorithm = "PS256"

	JWTIssuerEnvironmentKey            = "MONEYBAGS_JWT_ISSUER"
	JWTSigningAlgorithmEnvironmentKey  = "MONEYBAGS_JWT_SIGNING_ALGORITHM"
	JWTRSAPrivateKeyFileEnvironmentKey = "MONEYBAGS_JWT_RSA_PRIVATE_KEY_FILE"
)

type ContextKey string

const (
	UsernameContextKey ContextKey = "username"
)

var (
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserDoesNotExist = errors.New("user does not exist")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidEmail     = errors.New("invalid email")
)
