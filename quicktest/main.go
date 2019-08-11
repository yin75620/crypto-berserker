package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/yin75620/crypto-berserker/ftx"
)

type LoginRequest struct {
	Op   string             `json:"op"`
	Args LoginRequestDetail `json:"args"`
}

type LoginRequestDetail struct {
	Key        string `json:"key"`
	Time       string `json:"time"`       // integer current timestamp (in milliseconds)
	Sign       string `json:"sign"`       //SHA256 HMAC of the following string, using your API secret: <time>websocket_login
	Subaccount string `json:"subaccount"` // (optional) subaccount name
}

func main() {

	res, _ := ftx.GetParamHmacSHA256HexSign("Y2QTHI23f23f23jfjas23f23To0RfUwX3H42fvN-", "1557246346499"+ftx.WEBSOCKET_LOGIN_KEY_WORD)
	log.Println(res)
	loginRequest := LoginRequest{
		Op: "login",
		Args: LoginRequestDetail{
			Key:  "rrrr",
			Time: "aaaaa",
			Sign: "ooooo",
		},
	}

	loginJson, err := json.Marshal(loginRequest)
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(string(loginJson))
}
