package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	bsk "github.com/yin75620/crypto-berserker/setting"
)

var (
	apiURL    = "https://ftx.com/api"
	apiPrefix = "/api/"
)

func main() {
	client := http.DefaultClient
	method := "GET"
	path := "account"
	body := ""
	str := fmt.Sprintf("%s/%s", apiURL, path)
	req, err := http.NewRequest(method, str, nil)
	if err != nil {
		log.Println(err)
		return
	}

	addHeader(&req.Header, method, path, body)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("%s", err)
		return
	}

	fmt.Printf("%s", sitemap)
}

func addHeader(header *http.Header, reqMethod, path, body string) {
	nanos := time.Now().UnixNano() / 1000000
	ts := strconv.FormatInt(nanos, 10)

	header.Add("FTX-KEY", bsk.FTX_KEY)
	header.Add("FTX-TS", ts)
	//boyd := "" //之後再用
	payload := fmt.Sprintf("%s%s%s%s", ts, reqMethod, apiPrefix+path, body)
	log.Println(payload)
	sign, _ := GetParamHmacSHA256HexSign(bsk.FTX_API_SECRET_KEY, payload)
	log.Println(sign)
	header.Add("FTX-SIGN", sign)
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
