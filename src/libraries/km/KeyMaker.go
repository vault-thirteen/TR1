package km

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	WebTokenField_UserId         = "userId"
	WebTokenField_SessionId      = "sessionId"
	WebTokenField_ExpirationTime = "exp"
	TokenHeader_Alg              = "alg"
	TokenAlg_PS512               = "PS512" // RSA-PSS.
	TokenAlg_RS512               = "RS512" // RSA.
)

const (
	ErrSigningMethodIsNotSupported = "signing method is not supported"
	ErrTokenIsNotValid             = "token is not valid"
	ErrTokenIsBroken               = "token is broken"
	ErrFUnsupportedSigningMethod   = "unsupported signing method: %s"
	ErrFUnexpectedSigningMethod    = "unexpected signing method: %v"
	ErrTypeCast                    = "type cast error"
)

type KeyMaker struct {
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
	signingMethod     jwt.SigningMethod
	signingMethodName string
}

func New(signingMethodName string, privateKeyFilePath string, publicKeyFilePath string) (km *KeyMaker, err error) {
	var signingMethod jwt.SigningMethod
	switch signingMethodName {
	case TokenAlg_PS512:
		signingMethod = jwt.SigningMethodPS512
	case TokenAlg_RS512:
		signingMethod = jwt.SigningMethodRS512
	default:
		return nil, errors.New(ErrSigningMethodIsNotSupported)
	}

	km = &KeyMaker{
		signingMethod:     signingMethod,
		signingMethodName: strings.ToUpper(signingMethodName),
	}

	km.privateKey, err = getPrivateKey(privateKeyFilePath, signingMethodName)
	if err != nil {
		return nil, err
	}

	km.publicKey, err = getPublicKey(publicKeyFilePath, signingMethodName)
	if err != nil {
		return nil, err
	}

	return km, nil
}

func (km *KeyMaker) MakeJWToken(userId int, sessionId int, expirationTime time.Time) (tokenString string, err error) {
	claims := jwt.MapClaims{
		WebTokenField_UserId:         userId,
		WebTokenField_SessionId:      sessionId,
		WebTokenField_ExpirationTime: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(km.signingMethod, claims, nil)

	var s string
	s, err = token.SignedString(km.privateKey)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (km *KeyMaker) ValidateToken(tokenString string) (userId int, sessionId int, err error) {
	validator := NewValidator(km.signingMethodName, km.publicKey)

	var token *jwt.Token
	token, err = jwt.Parse(tokenString, validator.KeyFn)
	if err != nil {
		return validator.userId, validator.sessionId, err
	}

	if !token.Valid {
		return validator.userId, validator.sessionId, errors.New(ErrTokenIsNotValid)
	}

	return validator.userId, validator.sessionId, nil
}
