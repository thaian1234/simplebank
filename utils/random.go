package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopgrstuvwxyz"

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + r.Int63n(max-min+1) //min -> max
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := make([]string, 0, len(SupportedCurrencies))
	for k := range SupportedCurrencies {
		currencies = append(currencies, k)
	}
	n := len(currencies)
	return currencies[r.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email", RandomString(6))
}
