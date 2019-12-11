package keysmanager

import (
	"bytes"
	"fmt"
)

// Keys contains the *decrypted* keys
type Keys struct {
	Encryption, MAC []byte
}

// String prevents someone printing keys without realizing that they secret. If
// he *really* wants to see the keys, he has to print them manually
// (fmt.Println(keys.MAC))
func (k Keys) String() string {
	return fmt.Sprintf("[secret!] Keys{}")
}

// Equal compares whether the fields are equal
func (k Keys) Equal(target Keys) bool {
	return bytes.Equal(k.MAC, target.MAC) && bytes.Equal(k.Encryption, target.Encryption)
}
