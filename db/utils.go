package db

import (
	"path"
	"strings"
)

// JoinRootPath joins 2 paths, but guarantees that it will not go higher than
// root
func JoinRootPath(root, unsafePath string) string {
	unsafePath = path.Clean(unsafePath)
	if unsafePath[0] != '/' {
		unsafePath = "/" + unsafePath
	}
	for strings.HasPrefix(unsafePath, "/..") {
		unsafePath = unsafePath[3:]
	}
	return path.Join(root, unsafePath)
}
