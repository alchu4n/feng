package macho

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	// path := "/Applications/ClashX.app/Contents/MacOS/ClashX"
	// path := "/usr/bin/python3"
	path := "/Users/whoami/go/bin/ipa-medit"
	info, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range info {
		fmt.Printf("%#v", v)
	}
}
