package appvuln

import (
	"fmt"
	"testing"
)

func TestScanOne(t *testing.T) {
	t.Log(New().ScanSingleApp("/Applications/WeChat.app"))
}

func TestScan(t *testing.T) {
	vulns, err := New().Scan()
	if err != nil {
		t.Fatal(err)
	}

	for i := range vulns {
		fmt.Println(vulns[i])
	}
}
