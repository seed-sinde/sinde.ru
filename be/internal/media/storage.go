package media
import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)
type StorageObject struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
}
var (
	minioClientOnce sync.Once
	minioClient     *minio.Client
	minioClientErr  error
)
func MinIOEnabled() bool {
	return strings.TrimSpace(os.Getenv("MINIO_ENDPOINT")) != ""
}
func MinIOBucket() string {
	if raw := strings.TrimSpace(os.Getenv("MINIO_MEDIA_BUCKET")); raw != "" {
		return raw
	}
	return "media"
}
func ContentLengthHeader(size int64) string {
	if size <= 0 {
		return ""
	}
	return strconv.FormatInt(size, 10)
}
func minioEndpoint() (string, bool, error) {
	raw := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	if raw == "" {
		return "", false, errors.New("minio endpoint is not configured")
	}
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", false, err
		}
		if parsed.Host == "" {
			return "", false, errors.New("minio endpoint host is empty")
		}
		return parsed.Host, parsed.Scheme == "https", nil
	}
	return raw, false, nil
}
func storageClient() (*minio.Client, error) {
	minioClientOnce.Do(func() {
		endpoint, secure, err := minioEndpoint()
		if err != nil {
			minioClientErr = err
			return
		}
		minioClient, minioClientErr = minio.New(endpoint, &minio.Options{
			Creds: credentials.NewStaticV4(
				strings.TrimSpace(os.Getenv("MINIO_ROOT_USER")),
				strings.TrimSpace(os.Getenv("MINIO_ROOT_PASSWORD")),
				"",
			),
			Secure: secure,
		})
	})
	return minioClient, minioClientErr
}
func ensureMinIOBucket(ctx context.Context) error {
	client, err := storageClient()
	if err != nil {
		return err
	}
	bucket := MinIOBucket()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}
func detectContentType(storageKey string, data []byte) string {
	switch strings.ToLower(strings.TrimSpace(filepath.Ext(storageKey))) {
	case ".webp":
		return "image/webp"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".avif":
		return "image/avif"
	case ".svg":
		return "image/svg+xml"
	case ".glb":
		return "model/gltf-binary"
	case ".json":
		return "application/json"
	}
	if len(data) > 0 {
		return http.DetectContentType(data)
	}
	return "application/octet-stream"
}
func PutObjectBytes(ctx context.Context, storageKey string, data []byte, contentType string) error {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return errors.New("invalid media key")
	}
	if !MinIOEnabled() {
		path, err := AbsPath(key)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		return os.WriteFile(path, data, 0o644)
	}
	if err := ensureMinIOBucket(ctx); err != nil {
		return err
	}
	client, err := storageClient()
	if err != nil {
		return err
	}
	if strings.TrimSpace(contentType) == "" {
		contentType = detectContentType(key, data)
	}
	_, err = client.PutObject(ctx, MinIOBucket(), key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}
func OpenObject(ctx context.Context, storageKey string) (*StorageObject, error) {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return nil, errors.New("invalid media key")
	}
	if !MinIOEnabled() {
		path, err := ResolveFilePath(key)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		return &StorageObject{
			Body:        file,
			ContentType: detectContentType(key, nil),
			Size:        stat.Size(),
		}, nil
	}
	client, err := storageClient()
	if err != nil {
		return nil, err
	}
	object, err := client.GetObject(ctx, MinIOBucket(), key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	info, err := object.Stat()
	if err != nil {
		_ = object.Close()
		return nil, err
	}
	return &StorageObject{
		Body:        object,
		ContentType: strings.TrimSpace(info.ContentType),
		Size:        info.Size,
	}, nil
}
func ReadObjectBytes(ctx context.Context, storageKey string) ([]byte, error) {
	object, err := OpenObject(ctx, storageKey)
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()
	return io.ReadAll(object.Body)
}
func ObjectExists(ctx context.Context, storageKey string) (bool, error) {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return false, nil
	}
	if !MinIOEnabled() {
		path, err := ResolveFilePath(key)
		if err != nil {
			return false, err
		}
		_, err = os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	client, err := storageClient()
	if err != nil {
		return false, err
	}
	_, err = client.StatObject(ctx, MinIOBucket(), key, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}
	var response minio.ErrorResponse
	if errors.As(err, &response) && response.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, err
}
func DeleteObject(ctx context.Context, storageKey string) error {
	key := NormalizeStorageKey(storageKey)
	if key == "" {
		return nil
	}
	if !MinIOEnabled() {
		return removeFileAndEmptyParents(key)
	}
	client, err := storageClient()
	if err != nil {
		return err
	}
	return client.RemoveObject(ctx, MinIOBucket(), key, minio.RemoveObjectOptions{})
}
