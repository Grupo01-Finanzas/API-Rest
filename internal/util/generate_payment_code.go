package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

// GeneratePaymentCode generates a random 6-digit payment code
func GeneratePaymentCode() string {
	code := ""
	for i := 0; i < 6; i++ {
		// Generate a random number between 0 and 9
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// Handle the error appropriately (e.g., log and panic, or return an error)
			panic(fmt.Sprintf("error generating random number: %v", err))
		}
		code += strconv.Itoa(int(n.Int64()))
	}
	return code
}
