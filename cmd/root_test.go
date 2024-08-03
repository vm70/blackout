package cmd

import (
	"testing"
)

func TestRunApp(_ *testing.T) {
	for i := 0; i < 100; i++ {
		run(rootCmd, []string{"a very long message"})
	}
}
