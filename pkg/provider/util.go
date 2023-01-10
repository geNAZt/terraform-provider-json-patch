package provider

import (
	"crypto/sha1"
	"fmt"
)

func MakeId(s []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(s))
}
