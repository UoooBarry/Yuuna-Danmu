package live

import (
	"testing"
)

func TestAuth_FullCookie(t *testing.T) {
	session := &Session{}
	cookie := "DedeUserID=123456; buvid3=ABC-DEF-123;"

	err := session.prepareAuth(cookie)
	if err != nil {
		t.Fatalf("Error preparing auth: %v", err)
	}
	if session.authContext.UID != 123456 {
		t.Errorf("UID error: %d", session.authContext.UID)
	}
	if session.authContext.Buvid3 != "ABC-DEF-123" {
		t.Errorf("Buvid3 Error: %s", session.authContext.Buvid3)
	}
}

func TestAuth_GuestMode(t *testing.T) {
	session := &Session{}
	cookie := ""

	err := session.prepareAuth(cookie)
	if err != nil {
		t.Fatalf("Guest should not report error: %v", err)
	}
	if session.authContext.UID != 0 {
		t.Errorf("Guest should have UID 0, got %d", session.authContext.UID)
	}
	if session.authContext.Buvid3 == "" {
		t.Error("Guest missing Buvid3")
	}
}

func TestAuth_MissingBuvidError(t *testing.T) {
	session := &Session{}
	cookie := "DedeUserID=999999;"

	err := session.prepareAuth(cookie)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
