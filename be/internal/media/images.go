package media

import (
	"bytes"
	"errors"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"sync"
)

const (
	storageImageExt         = ".webp"
	storageImageWebPQuality = 92
	storageImageWebPMethod  = 6
)

var initImageCodecsOnce sync.Once

func initImageCodecs() {
	initImageCodecsOnce.Do(func() {
		avif.InitDecoder()
		webp.Init()
	})
}
func detectImageMIME(data []byte) string {
	sniffLen := 512
	if len(data) < sniffLen {
		sniffLen = len(data)
	}
	return http.DetectContentType(data[:sniffLen])
}
func decodeImage(data []byte, mime string) (image.Image, error) {
	initImageCodecs()
	reader := bytes.NewReader(data)
	switch mime {
	case "image/jpeg":
		return jpeg.Decode(reader)
	case "image/png":
		return png.Decode(reader)
	case "image/gif":
		return gif.Decode(reader)
	case "image/webp":
		return webp.Decode(reader)
	case "image/avif":
		return avif.Decode(reader)
	default:
		return nil, errors.New("unsupported image format")
	}
}
func normalizeImageForStorage(data []byte) ([]byte, string, error) {
	initImageCodecs()
	mime := detectImageMIME(data)
	img, err := decodeImage(data, mime)
	if err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, webp.Options{
		Quality: storageImageWebPQuality,
		Method:  storageImageWebPMethod,
		Exact:   true,
	}); err != nil {
		return nil, "", err
	}
	if buf.Len() == 0 {
		return nil, "", errors.New("image encoding failed")
	}
	return buf.Bytes(), storageImageExt, nil
}
func NormalizeImageForStorage(data []byte) ([]byte, string, error) {
	return normalizeImageForStorage(data)
}
