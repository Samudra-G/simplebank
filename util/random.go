package util

import (
	"math/rand"
	"strings"
)
const alphabet = "abcdefghijklmnopqrstuvwxyz"

//RandomInt generates a random int bewteen min and max
func RandomInt(min, max int64) int64{
	return min + rand.Int63n(max - min + 1)
}

//RandomString generates a random string of len n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandomOwner generates a new owner name
func RandomOwner() string {
	return RandomString(6)
}

//RandomMoney generates a random amt of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency generates a random currency
func RandomCurrency() string {
	currencies := []string {"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}