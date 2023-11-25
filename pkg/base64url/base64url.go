package base64url

import (
	"encoding/base64"
	"errors"
	"strings"
)

func Encode(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	str = strings.Replace(str, "+", "-", -1)
	str = strings.Replace(str, "/", "_", -1)
	str = strings.Replace(str, "=", "", -1)
	return str
}

func Decode(str string) ([]byte, error) {
	if strings.ContainsAny(str, "+/") {
		return nil, errors.New("invalid base64url encoding")
	}
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	for len(str)%4 != 0 {
		str += "="
	}
	return base64.StdEncoding.DecodeString(str)
}

func Base64RemovePadding(data string) []byte {
	bytes := []byte(data)
	unpadded := bytes
	for len(unpadded) > 0 && unpadded[len(unpadded)-1] == '=' {
		unpadded = unpadded[:len(unpadded)-1]
	}
	return unpadded
}
