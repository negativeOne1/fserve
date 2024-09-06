package signature

import (
	"encoding/hex"
	"errors"
)

type Validator interface {
	IsValid(algo, date, expires, method, resource, signature string) error
}

type HMACValidator struct {
	Secret string
}

func NewHMACValidator(secret string) *HMACValidator {
	return &HMACValidator{Secret: secret}
}

func (v *HMACValidator) IsValid(algo, date, expires, method, resource, signature string) error {
	s, err := CreateSignature(v.Secret, algo, date, expires, method, resource)
	if err != nil {
		return err
	}

	if hex.EncodeToString(s) != signature {
		return errors.New("Invalid signature")
	}

	return nil
}
