package maicoin

import (
	"testing"
	"time"
)

func TestWsStart(t *testing.T) {

	Start()

	time.Sleep(time.Duration(10) * time.Second)
}
