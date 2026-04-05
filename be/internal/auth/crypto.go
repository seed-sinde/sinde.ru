package auth
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)
func encryptString(plaintext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	out := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.RawStdEncoding.EncodeToString(out)
	return []byte(encoded), nil
}
func decryptString(ciphertext []byte, key []byte) (string, error) {
	raw, err := base64.RawStdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := raw[:gcm.NonceSize()]
	plain, err := gcm.Open(nil, nonce, raw[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
