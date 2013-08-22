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
	var result []byte

        for i := 0; i < len(pass); i++ {
		PPass, PKey := pass[i], ticket[i]
		APass, AKey := PPass >> 4, PKey % 16

		ANB, ANB2 := (APass + PKey) % uint8(len(alphanum)), (AKey + PKey) % uint8(len(alphanum))

		result = append(result, alphanum[ANB], alphanum[ANB2])
        }

	return string(result)
}

func DecryptDofusPassword(pass, ticket string) string {
	var PKey rune
	var ANB int
	var ANB2 int
	var somme1 int
	var somme2 int
	var APass int
	var AKey int

	var decrypted []byte

        for i := 0; i < len(pass); i += 2 {
		PKey = rune(ticket[i/2])
		ANB = strings.IndexRune(alphanum, rune(pass[i]))
		ANB2 = strings.IndexRune(alphanum, rune(pass[i+1]))

		somme1 = ANB + len(alphanum)
		somme2 = ANB2 + len(alphanum)

		APass = somme1 - int(PKey)
		if APass < 0 { APass += 64 }
		APass <<= 4

		AKey = somme2 - int(PKey)
		if AKey < 0 { AKey += 64 }

		PPass := byte(APass + AKey)

		decrypted = append(decrypted, PPass)
	}

	return string(decrypted)
}
