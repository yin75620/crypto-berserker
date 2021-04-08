package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	postWithJson()
}

const SOLONA_API = "https://solana-api.projectserum.com/"
const JsonMode = "application/json"

func postWithJson() {
	//post請求提交json數據
	sendContent := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTokenSupply",
		"params":  []interface{}{"4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R"},
	}

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
