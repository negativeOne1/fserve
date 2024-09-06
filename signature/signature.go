package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

var (
	SigningAlgorithm = "HMAC"
	HashingAlgorithm = "SHA256"
	AlgHashSeparator = "-"
	DefaultAlgorithm = SigningAlgorithm + AlgHashSeparator + HashingAlgorithm
)

func CreateSignature(algo, date, expires, method, resource string) ([]byte, error) {
	if algo != DefaultAlgorithm {
		return nil, fmt.Errorf("unsupported algorithm: %s", algo)
	}

	mac := hmac.New(sha256.New, []byte("secret"))
	fmt.Fprintf(mac, "%s:%s:%s:%s", date, expires, method, resource)

	return mac.Sum(nil), nil
}
