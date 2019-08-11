package ftx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

const (
	WEBSOCKET_LOGIN_KEY_WORD = "websocket_login"
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
