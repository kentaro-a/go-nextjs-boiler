package util

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

func MapToStruct(m interface{}, v interface{}) error {
	if b, err := json.Marshal(m); err != nil {
		return errors.WithStack(err)
	} else {

		if err := json.Unmarshal(b, v); err != nil {
			return errors.WithStack(err)
		}

	}
	return nil
}

func MakeRandStr(length int) string {
	if length > 62 {
		length = 62
	}
	base_chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	for i := range base_chars {
		j := rand.Intn(i + 1)
		base_chars[i], base_chars[j] = base_chars[j], base_chars[i]
	}
	return string(base_chars[:length])
}
