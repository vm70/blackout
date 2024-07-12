package main

import (
	"runtime"
	"testing"
)

func TestDownloadingPoems(t *testing.T) {
	downloadPoems("poems.json")
}


func TestPoemsLocation(t *testing.T) {
	if (runtime.GOOS[0:4] == "linux") && (dataFolder != "~/.local/share/blackout/") {
		t.Fail()
	}
}
