package message_tool

import (
	"testing"
)

func TestMain(t *testing.T) {
	StartTelegram()
	SendBroadcastArcherGroup("Hello I'm testing")
}
