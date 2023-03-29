package components

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	hashPre = "&*h1#$mkl)(23bh^&"
)

func CheckLogin(mobile string, password string) error {

	return nil
}

func EncryptPassword(password string) string {
	tool := md5.New()
	tool.Write([]byte(hashPre + password))
	return hex.EncodeToString(tool.Sum(nil))
}

func Login(input string, password string) bool {
	tool := md5.New()
	tool.Write([]byte(hashPre + input))
	input = hex.EncodeToString(tool.Sum(nil))
	return input == password
}
