package util

import (
	"math/rand"
	"strconv"
)

func GenerateVerificationCode() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
