package helpers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

func IntToBase64(num int) (string, error) {
	// Handle negative numbers if needed
	if num < 0 {
		return "", fmt.Errorf("negative numbers are not supported")
	}

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(num))

	// Trim leading zeros to make the output more compact
	// Find first non-zero byte
	start := 0
	for i := 0; i < len(bytes); i++ {
		if bytes[i] != 0 {
			start = i
			break
		}
	}

	return base64.StdEncoding.EncodeToString(bytes[start:]), nil
}

func Base64ToInt(base64Str string) (int, error) {
	// Decode Base64 string to bytes
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return 0, fmt.Errorf("invalid base64 string: %w", err)
	}

	// Pad the byte slice with leading zeros if necessary
	paddedBytes := make([]byte, 8)
	copy(paddedBytes[8-len(bytes):], bytes)

	// Convert bytes back to integer
	num := binary.BigEndian.Uint64(paddedBytes)

	// Check for integer overflow
	if num > uint64(^uint(0)>>1) {
		return 0, fmt.Errorf("number exceeds maximum int value")
	}

	return int(num), nil
}
