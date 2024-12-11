package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"platform/log"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func HexToBytes(dataHex string) ([]byte, error) {
	if len(dataHex)%2 != 0 {
		return []byte{}, fmt.Errorf("invalid hex string length")
	}
	data := make([]byte, len(dataHex)/2)
	_, err := hex.Decode(data, []byte(dataHex))
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func BytesToHex(data []byte) string {
	dataHex := make([]byte, len(data)*2)
	hex.Encode(dataHex, data)
	return string(dataHex)
}

func GetRand(size int) ([]byte, string, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.Errorf("Error generating %d random bytes %v", size, err)
		return []byte{}, "", err
	}

	return data, BytesToHex(data), nil
}

func HashPassword(password string, salt []byte) string {
	secret := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	return BytesToHex(secret)
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

func ParseTime(timeStr *string) (time.Time, error) {
	if timeStr == nil {
		return time.Time{}, fmt.Errorf("time string is nil")
	}
	s := strings.ReplaceAll(*timeStr, "T", " ")
	s = strings.ReplaceAll(s, "Z", "")
	s += strings.Repeat("0", 26-len(s))
	return time.Parse("2006-01-02 15:04:05.000000", s)
}
