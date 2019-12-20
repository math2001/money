// +build gofuzz

package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

	if strings.Contains(final, "/../") || strings.HasSuffix(final, "/..") {
		panic(fmt.Sprintf("path not safe %q %q in %q", root, unsafePath, final))
	}

	return 1
}
