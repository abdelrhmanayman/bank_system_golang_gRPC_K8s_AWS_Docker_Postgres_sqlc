package util

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func generateRandomString(stringLength int) string {
	b := make([]rune, stringLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateRandomOwner() string {
	return generateRandomString(8)
}

func CreateRandomBalance() int64 {
	return generateRandomInt(250, 25000)
}

func CreateRandomCurrency() string {
	currencies := []string{"USD", "EUR", "EGP"}
	currenciesLength := len(currencies)

	return currencies[rand.Intn(currenciesLength)]
}
