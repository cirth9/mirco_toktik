package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encoding 加密
func Md5Encoding(data string) string {
	sum := md5.Sum([]byte(data))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

// EncryptPassword 加密密码
func EncryptPassword(password string) string {
	return Md5Encoding(password)
}

// DecryptPassword 校验密码
func DecryptPassword(password string, encryptedPassword string) bool {
	return Md5Encoding(password) == encryptedPassword
}
