package handler

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	resultTestCode := m.Run()
	os.Exit(resultTestCode)
}
