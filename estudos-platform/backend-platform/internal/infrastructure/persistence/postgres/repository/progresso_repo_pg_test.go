package repository

import "testing"

func TestBoolToPercent(t *testing.T) {
	if boolToPercent(true) != 100 {
		t.Fatal("lido deve ser 100")
	}
	if boolToPercent(false) != 0 {
		t.Fatal("não lido deve ser 0")
	}
}
