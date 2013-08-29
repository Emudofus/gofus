package shared

import (
	"strings"
)

const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func Split2(in, sep string) (string, string) {
	res := strings.SplitN(in, sep, 2)
	return res[0], res[1]
}

func CryptDofusPassword(pass, ticket string) string {
	result := make([]byte, len(pass)*2)
	for i := 0; i < len(pass); i++ {
		PPass, PKey := int(pass[i]), int(ticket[i])
		APass, AKey := PPass>>4, PPass%16

		result[i*2] = alphanum[(APass+PKey)%len(alphanum)]
		result[i*2+1] = alphanum[(AKey+PKey)%len(alphanum)]
	}
	return string(result)
}

func DecryptDofusPassword(pass, ticket string) string {
	result := make([]byte, len(pass)/2)
	for i := 0; i < len(pass); i += 2 {
		PKey := int(ticket[i/2])
		ANB := strings.IndexRune(alphanum, rune(pass[i]))
		ANB2 := strings.IndexRune(alphanum, rune(pass[i+1]))

		somme1, somme2 := ANB+len(alphanum), ANB2+len(alphanum)

		APass, AKey := somme1-PKey, somme2-PKey

		if APass < 0 {
			APass += 64
		}
		APass <<= 4

		if AKey < 0 {
			AKey += 64
		}

		result[i/2] = byte(APass + AKey)
	}
	return string(result)
}
