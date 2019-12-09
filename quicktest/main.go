package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yin75620/crypto-berserker/jmath"
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
	arrayPrintTest()
}

func arrayPrintTest() {
	array := []interface{}{"test", "ABC"}
	//array := []string{"123", "456"}
	fmt.Println(fmt.Sprintf("TEST:%v", array))
}

func indexTest() {
	s := "abcdjefghij"
	lastIndex := strings.LastIndex(s, "j")
	fmt.Println(lastIndex)

}

func gorountineTest() {
	const MAX = 3
	countChannel := make(chan int)
	for i := 0; i < MAX; i++ {
		go func() {
			d := time.Second * time.Duration(i)
			time.Sleep(d)
			fmt.Println("Do samething")
			countChannel <- 0
		}()
	}

	for i := 0; i < MAX; i++ {
		fmt.Println("ok")
		<-countChannel
	}

	fmt.Println("finish")
}

func dotTest() {
	askVolumeStr := strToFloat64(fmt.Sprintf("%.3f", 2.000690), 10)

	fmt.Println(askVolumeStr)

	fmt.Println(jmath.FloatFloorByFloat(1.3251, 0))
	//fmt.Println(jmath.FloatFloorByFloat(1.2541876, 0.0001))
	fmt.Println(jmath.FloatFloorByFloat(1.3251876110011, 0.0001))
}

func FloatFloorByFloat(f float64, unit float64) float64 {
	if unit == 0 {
		return f //不運算
	}
	val := f / unit
	val = math.Floor(val)
	res := val * unit / math.Pow10(-int(math.Log10(unit))) //為了不要出現浮點數
	return res
}

func FloatFloor(f float64, count int) float64 {
	val := f * math.Pow10(count)
	val = math.Floor(val)
	res := val / math.Pow10(count)
	return res
}

func strToFloat64(str string, len int) float64 {
	lenstr := "%." + strconv.Itoa(len) + "f"
	value, _ := strconv.ParseFloat(str, 64)
	value = math.Floor(math.Pow10(len)*value) / math.Pow10(len) // 無條件捨去
	nstr := fmt.Sprintf(lenstr, value)
	val, _ := strconv.ParseFloat(nstr, 64)
	return val
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
