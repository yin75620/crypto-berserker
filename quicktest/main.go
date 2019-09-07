package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
	timeTest()
}

func timeTest() {
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	fmt.Println(fileName)
}

func nameTest() {
	fmt.Println(os.Args[0])
	fmt.Println(os.Args[1:])
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
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
