package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	USD      = "USD"
	EURO     = "EURO"
	CAD      = "CAD"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandOwner() string {
	return RandString(6)
}

func RandBalance() int64 {
	return RandomInt(100, 100000)
}

func RandomCurrency() string {
	currencies := []string{USD, EURO, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
