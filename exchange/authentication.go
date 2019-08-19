package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"time"
)

func GetTimeSpanStr(nanos int64) string {
	ts := strconv.FormatInt(nanos, 10)
	return ts
}

func GetTimeSpan() int64 {
	nanos := time.Now().UnixNano() / 1000000
	return nanos
}

func GetParamHmacSHA256HexSign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return hex.EncodeToString(signByte), nil
}

func GetParamHmacSHA256Base64Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signByte), nil
}
