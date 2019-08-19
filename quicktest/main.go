package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
	timerTest()
}

func codeText() {
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

func orderTest() {

	var ftxClient = ftx.NewFtx(http.DefaultClient)

	var myOrder ftx.FtxOrder = ftx.FtxOrder{
		Market:    "FTT/USD",
		Side:      "sell",
		Price:     1.70,
		Size:      1,
		OrderType: ftx.MARKET,
	}
	response := ftxClient.PostOrder(myOrder)
	fmt.Println(response)
}

func logTest() {

	/*f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "prefix", log.LstdFlags)
	logger.Println("text to append")
	logger.Println("more text to append")*/

	logFile, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//log.SetOutput(f)
	log.Println("This is a test log entry")
}

func timerTest() {
	var delay_time int = 5
	d := time.Duration(time.Second * time.Duration(delay_time))

	t := time.NewTimer(d)
	defer t.Stop()

	var count = 0
	for {
		<-t.C
		plusSecond := 0
		if count > 2 {
			plusSecond = -4
		}
		count = count + 1
		t.Reset(time.Second * time.Duration(delay_time+plusSecond))
		time.Sleep(time.Second * 3)

		log.Println("TEST")
	}
}
