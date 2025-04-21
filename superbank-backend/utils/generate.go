package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	countryCode   = "IDN"
	branchCode    = "000001"
	identifierLen = 4
	accountLen    = 8
)

func generateRandomString(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomDigits(n int) string {
	max := 1
	for i := 0; i < n; i++ {
		max *= 10
	}
	return fmt.Sprintf("%0*d", n, rand.Intn(max))
}

func calculateCheckDigit(input string) string {
	sum := 0
	for _, c := range input {
		sum += int(c)
	}
	return fmt.Sprintf("%02d", sum%100)
}

func GenerateBankAccountNumber() string {
	rand.Seed(time.Now().UnixNano())

	identifier := generateRandomString(identifierLen)
	accountNumber := generateRandomDigits(accountLen)

	base := strings.Join([]string{identifier, branchCode, accountNumber}, "")
	checkDigit := calculateCheckDigit(base)

	return fmt.Sprintf("%s-%s-%s-%s-%s", countryCode, checkDigit, identifier, branchCode, accountNumber)
}
