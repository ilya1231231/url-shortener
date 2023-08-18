package random

import "math/rand"

func MakeRandStr(size int) string {
	charRunes := []rune("ABCDEFGHYJKLMNOPQRSTUVWXYZabcdefghyjklmnopqrstyvwxyz")
	b := make([]rune, size)
	for i := 0; i < size; i++ {
		b[i] = charRunes[rand.Intn(len(charRunes))]
	}
	return string(b)
}
