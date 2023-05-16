package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func Generate_unique_account_num() string {
	const account_num_len = 13

	rand.Seed(time.Now().UnixNano())

	set_num := "0123456789"
	set_num_len := len(set_num)

	result := make([]byte, account_num_len)
	for i := 0; i < account_num_len; i++ {
		random_index := rand.Intn(set_num_len)
		result[i] = set_num[random_index]
	}

	return string(result)
}

func Hash_password(_s_password string) string {
	hash := sha256.Sum256([]byte(_s_password))
	return hex.EncodeToString(hash[:])
}
