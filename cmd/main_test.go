package main

import (
	"os"
	"testing"
)

func TestMainCLI(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd"}
	main()

	os.Args = []string{"cmd", "distance", "foo"}
	main()

	os.Args = []string{"cmd", "distance", "foo", "bar"}
	main()

	os.Args = []string{"cmd", "closest", "foo"}
	main()

	os.Args = []string{"cmd", "closest", "foo", "bar", "baz"}
	main()

	os.Args = []string{"cmd", "unknown", "foo"}
	main()
}
