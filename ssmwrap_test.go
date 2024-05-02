package ssmwrap

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitLogger()
	os.Exit(m.Run())
}
