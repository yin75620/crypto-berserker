package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	getAccountInfo("56kGnD4FhM6uB1m1cZmNaaAxAiRJJTrYWCunCYANQkW1")
	getBalance("56kGnD4FhM6uB1m1cZmNaaAxAiRJJTrYWCunCYANQkW1")
}

const SOLONA_API = "https://solana-api.projectserum.com/"
const JsonMode = "application/json"

func getTokenSupply() {
	//post請求提交json數據
	sendContent := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenSupply",
		"params":  []interface{}{"4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"},
	}
	doPostWithJson(sendContent)
}

func getAccountInfo(pubkey string) {
	//post請求提交json數據
	sendContent := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getAccountInfo",
		"params":  []interface{}{pubkey},
	}
	doPostWithJson(sendContent)
}

func getBalance(pubkey string) {
	//post請求提交json數據
	sendContent := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params":  []interface{}{pubkey},
	}
	doPostWithJson(sendContent)
}

func doPostWithJson(sendContent map[string]interface{}) {
	jContent, err := json.Marshal(sendContent)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post(SOLONA_API, JsonMode, bytes.NewBuffer([]byte(jContent)))
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Post request with json result: %s\n", string(body))

}
