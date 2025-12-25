package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
)

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

func ParseActivityTime(raw string) (time.Time, error) {
	value := strings.TrimSpace(raw)
	value = strings.ReplaceAll(value, "T", " ")
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("时间格式不正确")
	}
	return time.Time{}, lastErr
}
