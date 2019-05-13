package sha256

import (
    "crypto/sha256"
    "encoding/hex"
)

func Sha256Sum(data []byte) string {
    bytes := sha256.Sum256(data)
    return hex.EncodeToString(bytes[:])
}
