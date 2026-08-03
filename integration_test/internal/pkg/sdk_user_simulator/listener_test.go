package sdk_user_simulator

import (
	"context"
	"testing"
)

func TestSendMsgCallBackListener_WaitReturnsAfterSuccess(t *testing.T) {
	listener := NewTestSendMsgCallBackListener("sender")

	listener.OnSuccess("")

	if err := listener.Wait(context.Background()); err != nil {
		t.Fatalf("Wait() error = %v, want nil", err)
	}
}

func TestSendMsgCallBackListener_WaitReturnsSendError(t *testing.T) {
	listener := NewTestSendMsgCallBackListener("sender")

	listener.OnError(1001, "send failed")

	if err := listener.Wait(context.Background()); err == nil {
		t.Fatal("Wait() error = nil, want send error")
	}
}
