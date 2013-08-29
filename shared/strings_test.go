package shared

import (
	"testing"
)

func TestCryptDofusPassword(t *testing.T) {
	password, ticket := "test", "azertyuiopqsdfghjklmwxcvbn012345"

	expected := "OLa_SO52"

	if test := CryptDofusPassword(password, ticket); test != expected {
		t.Errorf("%s is not equal to %s", test, expected)
	}
}

func TestDecryptDofusPassword(t *testing.T) {
	password, ticket := "test", "azertyuiopqsdfghjklmwxcvbn012345"

	test := DecryptDofusPassword(CryptDofusPassword(password, ticket), ticket)

	if test != password {
		t.Errorf("%s is not equal to %s", test, password)
	}
}
