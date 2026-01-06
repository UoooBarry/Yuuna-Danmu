package ui

import (
	"context"
	"testing"

	"uooobarry/yuuna-danmu/pkg/live"
)

type eventCaptured struct {
	name string
	data any
}

func TestWailsUI_Events(t *testing.T) {
	var captured *eventCaptured

	ui := NewWailsUI(nil)
	ui.emitter = func(ctx context.Context, eventName string, opts ...any) {
		captured = &eventCaptured{
			name: eventName,
			data: opts[0],
		}
	}

	t.Run("AppendSysMsg", func(t *testing.T) {
		captured = nil
		ui.SetContext(context.Background())

		msg := "系统已就绪"
		ui.AppendSysMsg(msg)

		if captured == nil || captured.name != live.SysMsgEvent {
			t.Errorf("Cannot capture sys:message event")
		}
		if captured.data != msg {
			t.Errorf("data no matched: %v", captured.data)
		}
	})

	t.Run("AppendDanmu validation", func(t *testing.T) {
		captured = nil
		ui.SetContext(context.Background())

		ui.AppendDanmu("Yuuna", 12, "花花人花花", "你被逮捕了")

		if captured == nil {
			t.Fatal("未能捕获弹幕事件")
		}

		data, ok := captured.data.(map[string]string)
		if !ok {
			t.Fatal("Dannmu type error")
		}

		expectedFields := map[string]string{
			"nickname":   "花花人花花",
			"content":    "你被逮捕了",
			"medalName":  "Yuuna",
			"medalLevel": "12",
		}

		for field, expected := range expectedFields {
			if data[field] != expected {
				t.Errorf("Validation %s failed: expect %s, got %s", field, expected, data[field])
			}
		}

		if data["timestamp"] == "" {
			t.Error("No timestamp")
		}
	})
}
