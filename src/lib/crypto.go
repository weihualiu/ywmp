package lib

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)


func SHA128_Buf(data []byte) string {
	//sha := sha1.New()
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func SHA128(key string) string {
	sha := sha1.New()
	io.WriteString(sha, key)
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func SHA256(key string) string {
	sha := sha256.New()
    io.WriteString(sha, key)
    return fmt.Sprintf("%x", sha.Sum(nil))
}

func SHA256_Buf(data []byte) string {
	//sha := sha256.New()
	return fmt.Sprintf("%x", sha256.Sum256(data))
}