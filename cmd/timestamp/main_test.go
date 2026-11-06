package main

import (
	"bytes"
	"testing"
)

func TestVersion(t *testing.T) {
	oldVersion, oldRevision := version, revision
	version, revision = "v1.2.3", "abcdef0"
	t.Cleanup(func() {
		version, revision = oldVersion, oldRevision
	})

	var out bytes.Buffer
	if err := run([]string{"-version"}, &out); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "timestamp v1.2.3 (rev: abcdef0)\n"; got != want {
		t.Errorf("version output = %q, want %q", got, want)
	}
}
