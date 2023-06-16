package signature

import (
	"crypto/md5"
	"encoding/hex"
)

func ConfuseHexMd5(k any) any {
	switch v := k.(type) {
	case []byte:
		val := md5.Sum(v)
		dst := make([]byte, hex.EncodedLen(len(val)))
		hex.Encode(dst, val[:])
		return dst
	case string:
		val := md5.Sum([]byte(v))
		return hex.EncodeToString(val[:])
	default:
		return k
	}
}
