package live

import (
	"testing"
)

func init() {
	session := &Session{}
	session.prepareAuth("")
}

func TestGetRealRoomID(t *testing.T) {
	rid, err := GetRealRoomID(23990839)
	if err != nil {
		t.Fatalf("GetRealRoomID failed: %v", err)
	}
	t.Logf("Real room ID: %d", rid)
}

func TestGetDanmuConfig(t *testing.T) {
	rid, err := GetRealRoomID(23990839)
	if err != nil {
		t.Fatalf("GetRealRoomID failed: %v", err)
	}
	_, err = GetDanmuConfig(rid)
	if err != nil {
		t.Fatalf("GetDanmuConfig failed: %v", err)
	}
}
