package shared

import (
	"math/rand"
	"time"
	"testing"
)

func TestCryptDofusPassword(t *testing.T) {
	password, ticket := "test", NextString(rand.NewSource(time.Now().UnixNano()), 32)
	expected := "TODO"

	if test := CryptDofusPassword(password, ticket); test != expected {
		t.Errorf("%s is not equal to %s", test, expected)
	}
}

func TestDecryptDofusPassword(t *testing.T) {
	password, ticket := "test", NextString(rand.NewSource(time.Now().UnixNano()), 32)

	test := DecryptDofusPassword(CryptDofusPassword(password, ticket), ticket)

	if test != password {
		t.Errorf("%s is not equal to %s", test, password)
	}
}
