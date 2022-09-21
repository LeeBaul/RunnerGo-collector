package server

import "testing"

func TestSendStopMsg(t *testing.T) {
	SendStopMsg("kpcontroller.apipost.cn:443", "700")
}
