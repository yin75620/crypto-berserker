package simpleLog

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func StartFileLog(fileName string) *os.File {

	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return logFile
}

func StartLog() *os.File {
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))

	return StartFileLog(fileName)
}
