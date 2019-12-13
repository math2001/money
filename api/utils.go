package api

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

func scryptKey(payload, salt []byte) []byte {
	const keysize = 32
	k, err := scrypt.Key(payload, salt, 1<<15, 8, 1, keysize)
	// the only possible errors returned by scrypt report about wrong
	// parameters (the numbers). ie. there is nothing a user of scryptKey can
	// do about an error it would get from it. Hence, we panic
	if err != nil {
		panic(fmt.Sprintf("wrong parameters for scrypt key: %s", err))
	}
	return k
}
