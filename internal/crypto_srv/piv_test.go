package crypto_srv

import (
	"testing"
)

func TestHelloName(t *testing.T) {
	err := NewYPIV("")

	if err == nil {
		t.Fail()
	}

}
