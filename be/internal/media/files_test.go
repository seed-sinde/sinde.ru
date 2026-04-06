package media

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"path/filepath"
	"testing"
)

func TestNormalizeStorageKey(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{raw: " users/avatars/test.webp ", want: "users/avatars/test.webp"},
		{raw: `users\avatars\test.webp`, want: "users/avatars/test.webp"},
		{raw: "/users/avatars/test.webp", want: "users/avatars/test.webp"},
		{raw: "../secret", want: ""},
		{raw: "users/../../secret", want: ""},
		{raw: "users/<bad>.webp", want: ""},
		{raw: "", want: ""},
	}
	for _, tc := range cases {
		if got := NormalizeStorageKey(tc.raw); got != tc.want {
			t.Fatalf("NormalizeStorageKey(%q): expected %q, got %q", tc.raw, tc.want, got)
		}
	}
}
func TestMaxBytesUsesDefaultAndValidOverride(t *testing.T) {
	t.Setenv("MEDIA_IMAGE_MAX_BYTES", "")
	if got := MaxBytes(); got != 8*1024*1024 {
		t.Fatalf("unexpected default max bytes: %d", got)
	}
	t.Setenv("MEDIA_IMAGE_MAX_BYTES", "12345")
	if got := MaxBytes(); got != 12345 {
		t.Fatalf("unexpected env max bytes: %d", got)
	}
	t.Setenv("MEDIA_IMAGE_MAX_BYTES", "-1")
	if got := MaxBytes(); got != 8*1024*1024 {
		t.Fatalf("expected default for invalid value, got %d", got)
	}
}
func TestAbsPathUsesRootDirAndRejectsInvalidKeys(t *testing.T) {
	root := t.TempDir()
	t.Setenv("MEDIA_FILES_DIR", root)
	path, err := AbsPath("users/avatars/test.webp")
	if err != nil {
		t.Fatalf("AbsPath returned error: %v", err)
	}
	want := filepath.Join(root, "users", "avatars", "test.webp")
	if path != want {
		t.Fatalf("expected path %q, got %q", want, path)
	}
	if _, err := AbsPath("../secret"); err == nil {
		t.Fatal("expected invalid path error for traversal key")
	}
}
func TestMediaProxyURLAndNormalizeImageForStorage(t *testing.T) {
	if got := MediaProxyURL(" users/avatars/test.webp "); got != "/api/proxy/media/files/users/avatars/test.webp" {
		t.Fatalf("unexpected media proxy url: %q", got)
	}
	if got := MediaProxyURL("../secret"); got != "" {
		t.Fatalf("expected empty url for invalid key, got %q", got)
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 1, color.RGBA{G: 255, A: 255})
	var src bytes.Buffer
	if err := png.Encode(&src, img); err != nil {
		t.Fatalf("failed to encode source png: %v", err)
	}
	normalized, ext, err := NormalizeImageForStorage(src.Bytes())
	if err != nil {
		t.Fatalf("NormalizeImageForStorage returned error: %v", err)
	}
	if ext != ".webp" {
		t.Fatalf("unexpected extension: %q", ext)
	}
	if len(normalized) == 0 {
		t.Fatal("expected normalized image bytes")
	}
	if got := detectImageMIME(normalized); got != "image/webp" {
		t.Fatalf("unexpected normalized mime: %q", got)
	}
}
