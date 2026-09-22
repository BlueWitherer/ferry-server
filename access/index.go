package access

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
)

func HashString(b []byte) (string, string) {
	raw := base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(raw))
	return raw, base64.RawURLEncoding.EncodeToString(h[:])
}

func GetDomain(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func FullURL(r *http.Request) string {
	base := GetDomain(r)
	return fmt.Sprintf("%s%s", base, r.RequestURI)
}
