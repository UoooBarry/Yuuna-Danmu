package wbi

import (
	"testing"
)

func TestRefreshBuvid3(t *testing.T) {
	err := RefreshBuvid3()
	if err != nil {
		t.Fatalf("RefreshBuvid3 failed: %v", err)
	}

	cookie, err := GetBuvid3()
	if err != nil {
		t.Fatalf("GetBuvid3 failed: %v", err)
	}
	if cookie == "" {
		t.Error("RefreshBuvid3 failed: buvid3 cookie not found")
	}
}
