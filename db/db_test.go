package db

import (
	"path/filepath"
	"testing"
)

const testdir = "test-tmp"

func TestIntegration(t *testing.T) {
	app, err := NewApp(filepath.Join(testdir, "test-"+t.Name()))
	if err != nil {
		t.Fatalf("making new app: %s", err)
		return
	}
	defer app.RemoveAllDataForever()

	// register a new user A
	// register a new user B
	// sign in as user A and compare that we get the same data
	// save info as user A
	// try to load data as user A
	// try to load the same data as user B
	// save same data as user B
	// make sure that user A info is still the original, and that user B info
	// is what we just put in
}
