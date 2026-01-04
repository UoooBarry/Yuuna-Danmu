package live

import (
	"testing"

	"uooobarry/yuuna-danmu/pkg/wbi"
)

func init() {
	_ = wbi.EnsureBuvid()
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
	rsp, err := GetDanmuConfig(rid)
	if err != nil {
		t.Fatalf("GetDanmuConfig failed: %v", err)
	}
	t.Logf("Danmu config: %+v", rsp)
}
