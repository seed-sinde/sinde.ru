package auth

import "testing"

func TestBuildDeviceContextFingerprintIgnoresIP(t *testing.T) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4)"
	acceptLanguage := "ru-RU,ru;q=0.9,en-US;q=0.8"
	first := BuildDeviceContext("10.20.30.40", userAgent, acceptLanguage)
	second := BuildDeviceContext("203.0.113.24", userAgent, acceptLanguage)
	if first.FingerprintHash != second.FingerprintHash {
		t.Fatalf("expected identical fingerprints for different IPs, got %q and %q", first.FingerprintHash, second.FingerprintHash)
	}
	if first.IP != "10.20.30.40" {
		t.Fatalf("expected first IP to be preserved, got %q", first.IP)
	}
	if second.IP != "203.0.113.24" {
		t.Fatalf("expected second IP to be preserved, got %q", second.IP)
	}
}
func TestBuildDeviceContextFingerprintNormalizesAcceptLanguage(t *testing.T) {
	userAgent := "Mozilla/5.0"
	first := BuildDeviceContext("203.0.113.10", userAgent, "ru-RU, ru;q=0.9, en-US;q=0.8")
	second := BuildDeviceContext("203.0.113.10", userAgent, "ru-ru,ru,en-us")
	if first.FingerprintHash != second.FingerprintHash {
		t.Fatalf("expected normalized Accept-Language fingerprints to match, got %q and %q", first.FingerprintHash, second.FingerprintHash)
	}
}
func TestBuildDeviceContextDeviceLabelSummarizesUserAgent(t *testing.T) {
	device := BuildDeviceContext(
		"203.0.113.10",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36",
		"ru-RU",
	)
	if device.DeviceLabel != "Chrome on Linux" {
		t.Fatalf("expected summarized device label, got %q", device.DeviceLabel)
	}
	if device.UserAgent == device.DeviceLabel {
		t.Fatalf("expected device label to differ from raw user agent")
	}
}
