package main

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var (
	secretKey       = "secret"
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
	Iat  string `json:"iat"`
}

// Implement function Valid in Claim interface
func (payload *Payload) Valid() error {
	return nil
}

func CreateToken() (string, error) {
	payload := Payload{
		Sub:  "1234567890",
		Name: "John Doe",
		Iat:  "1516239022",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	return jwtToken.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
