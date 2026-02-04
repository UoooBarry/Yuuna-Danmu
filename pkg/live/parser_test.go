package live

import (
	"testing"
)

func TestParseToast(t *testing.T) {
	jsonBody := `{
		"cmd": "USER_TOAST_MSG",
		"data": {
			"uid": 123456,
			"username": "TestUser",
			"guard_level": 3,
			"price": 138000,
			"num": 1,
			"unit": "month",
			"role_name": "舰长",
			"toast_msg": "TestUser just became a Captain!"
		}
	}`

	data := parseToast([]byte(jsonBody))

	if data == nil {
		t.Fatal("Expected non-nil data")
	}

	if data.UID != 123456 {
		t.Errorf("Expected UID 123456, got %d", data.UID)
	}
	if data.Username != "TestUser" {
		t.Errorf("Expected Username TestUser, got %s", data.Username)
	}
	if data.GuardLevel != Captain {
		t.Errorf("Expected GuardLevel %d (Captain), got %d", Captain, data.GuardLevel)
	}
	if data.Price != 138000 {
		t.Errorf("Expected Price 138000, got %d", data.Price)
	}
	if data.Num != 1 {
		t.Errorf("Expected Num 1, got %d", data.Num)
	}
	if data.Unit != "month" {
		t.Errorf("Expected Unit month, got %s", data.Unit)
	}
	if data.RoleName != "舰长" {
		t.Errorf("Expected RoleName 舰长, got %s", data.RoleName)
	}
}

func TestDispatchUserToast(t *testing.T) {
	// This test verifies that dispatch calls parseToast and sends event to channel
	eventCh := make(chan Event, 1)
	c := &WsClient{
		eventCh: eventCh,
	}

	jsonBody := `{
		"cmd": "USER_TOAST_MSG",
		"data": {
			"uid": 666,
			"username": "SuperFan",
			"guard_level": 1,
			"price": 2000000,
			"num": 12,
			"unit": "month",
			"role_name": "总督"
		}
	}`

	c.dispatch([]byte(jsonBody))

	select {
	case event := <-eventCh:
		if event.Type != UserToastEvent {
			t.Errorf("Expected event type %s, got %s", UserToastEvent, event.Type)
		}
		toastData, ok := event.Data.(*ToastMsgData)
		if !ok {
			t.Fatalf("Expected event data to be *ToastMsgData")
		}
		if toastData.UID != 666 {
			t.Errorf("Expected UID 666, got %d", toastData.UID)
		}
		if toastData.GuardLevel != Governor {
			t.Errorf("Expected GuardLevel %d (Governor), got %d", Governor, toastData.GuardLevel)
		}
	default:
		t.Fatal("Expected event in channel, but got none")
	}
}

func TestParseGiftStarProcess(t *testing.T) {
	jsonBody := `{
		"cmd": "GIFT_STAR_PROCESS",
		"data": {
			"message": "用户 TestUser 的礼物星程已更新"
		}
	}`

	data := parseGiftStarProcess([]byte(jsonBody))

	if data == nil {
		t.Fatal("Expected non-nil data")
	}

	if data.Message != "用户 TestUser 的礼物星程已更新" {
		t.Errorf("Expected message '用户 TestUser 的礼物星程已更新', got %s", data.Message)
	}
}

func TestDispatchGiftStarProcess(t *testing.T) {
	eventCh := make(chan Event, 1)
	c := &WsClient{
		eventCh: eventCh,
	}

	jsonBody := `{
		"cmd": "GIFT_STAR_PROCESS",
		"data": {
			"message": "礼物星程进度更新"
		}
	}`

	c.dispatch([]byte(jsonBody))

	select {
	case event := <-eventCh:
		if event.Type != GiftStarProcessEvent {
			t.Errorf("Expected event type %s, got %s", GiftStarProcessEvent, event.Type)
		}
		giftStarData, ok := event.Data.(*GiftStarProcessData)
		if !ok {
			t.Fatalf("Expected event data to be *GiftStarProcessData")
		}
		if giftStarData.Message != "礼物星程进度更新" {
			t.Errorf("Expected message '礼物星程进度更新', got %s", giftStarData.Message)
		}
	default:
		t.Fatal("Expected event in channel, but got none")
	}
}
