// +build gofuzz

package db

import (
	"bytes"
	"encoding/binary"
	"path"
	"strings"

	fuzz "github.com/google/gofuzz"
)

// this is very naive I'm guessing, but it fuzzes
func Fuzz(data []byte) int {
	if len(data) != 8 {
		return 0
	}

	var seed int64
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &seed)
	f := fuzz.NewWithSeed(seed)

	var root string
	var unsafePath string
	f.Fuzz(&root)
	f.Fuzz(&unsafePath)

	final := JoinRootPath(root, unsafePath)
	final = path.Clean(final)

	if !strings.HasPrefix(final, root) {
		panic("invalid prefix")
	}
	if !strings.Contains(final, "/../") {
		panic("path not safe")
	}

	return 1
}
