package tools

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var sb strings.Builder
	k := len(letters)

	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.Intn(k)])
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return RandomString(10) + "@" + RandomString(5) + ".com"
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 100000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
