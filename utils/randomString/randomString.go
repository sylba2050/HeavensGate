package randomString

import (
    "math/rand"
    "time"
)

const rs2Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandString(n int) string {
    rand.Seed(time.Now().UnixNano())

    b := make([]byte, n)
    for i := range b {
        b[i] = rs2Letters[rand.Intn(len(rs2Letters))]
    }

    return string(b)
}
