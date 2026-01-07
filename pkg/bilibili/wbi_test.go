package bilibili

import (
	"net/url"
	"testing"
)

var testRawURL string = "https://api.bilibili.com/x/space/wbi/acc/info?mid=688134497"

// TestWbiUpdate 测试获取/更新 Key
func TestWbiUpdate(t *testing.T) {
	err := Update()
	if err != nil {
		t.Fatalf("Update WBI keys failed: %v", err)
	}

	keys, err := Get()
	if err != nil {
		t.Fatalf("Get WBI keys failed: %v", err)
	}

	if wbiKeys.Img == "" || wbiKeys.Sub == "" {
		t.Error("Fetched keys are empty")
	}
	t.Logf("Successfully fetched keys: ImgKey=%s, SubKey=%s", keys.Img, keys.Sub)
}

// TestWbiSign 测试签名追加功能
func TestWbiSign(t *testing.T) {
	// 确保 keys 已初始化
	if _, err := Get(); err != nil {
		t.Skip("Skipping sign test because keys couldn't be fetched")
	}

	u, _ := url.Parse(testRawURL)

	err := Sign(u)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	query := u.Query()

	// 验证是否包含了必填的 WBI 签名参数
	if query.Get("w_rid") == "" {
		t.Error("Missing 'w_rid' in signed URL")
	}
	if query.Get("wts") == "" {
		t.Error("Missing 'wts' (timestamp) in signed URL")
	}

	t.Logf("Signed URL: %s", u.String())
}

// TestWbiConsistency 测试多次签名的一致性/变动性
func TestWbiConsistency(t *testing.T) {
	u, _ := url.Parse(testRawURL)

	Sign(u)
	rid1 := u.Query().Get("w_rid")

	// WBI 签名包含时间戳 wts，如果运行极快，wts 可能相同，w_rid 也应相同
	u2, _ := url.Parse(testRawURL)
	Sign(u2)
	rid2 := u2.Query().Get("w_rid")

	if rid1 == "" || rid2 == "" {
		t.Error("Generated empty w_rid")
	}
}
