package base58

import "github.com/jbenet/go-base58"

func BytesFromB58(b58 string) []byte {
	return base58.Decode(b58)
}

func BytesToB58(p []byte) string {
	return base58.Encode(p)
}
