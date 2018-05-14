package restapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func authorizationHeader(method, path string, now time.Time, apiKey, secret string) http.Header {
	date := now.Format(time.RFC1123Z)

	authorization := fmt.Sprintf("SUMO %s:%s", apiKey,
		computeHmacSha1(fmt.Sprintf("%s\n%s\n%s", method, path, date), secret))

	return http.Header{
		"Date":          {date},
		"Authorization": {authorization},
	}
}

func computeHmacSha1(message string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
