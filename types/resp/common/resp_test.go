package common

import "testing"

func TestResponse_SetNoData(t *testing.T) {
	r := &Response{}
	r.SetNoData(CodeSuccess)
	if r.StatusCode != CodeSuccess {
		t.Fatalf("code mismatch: %d", r.StatusCode)
	}
	if r.StatusMsg != Msg[CodeSuccess] {
		t.Fatalf("msg mismatch: %s", r.StatusMsg)
	}
	if r.Data != nil {
		t.Fatalf("expected nil data, got: %#v", r.Data)
	}
}

func TestResponse_SetWithData(t *testing.T) {
	r := &Response{}
	r.SetWithData(CodeInvalidParams, 123)
	if r.StatusCode != CodeInvalidParams {
		t.Fatalf("code mismatch: %d", r.StatusCode)
	}
	if r.StatusMsg != Msg[CodeInvalidParams] {
		t.Fatalf("msg mismatch: %s", r.StatusMsg)
	}
	if r.Data != 123 {
		t.Fatalf("data mismatch: %#v", r.Data)
	}
}

func TestGetMsg(t *testing.T) {
	if GetMsg(CodeSuccess) != "success" {
		t.Fatal("unexpected message for success")
	}
	if GetMsg(999999) != "" {
		t.Fatal("unknown code should return empty string")
	}
}
