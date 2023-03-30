package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	for _, tt := range []struct {
		p1     string
		p2     string
		result bool
	}{
		{p1: "123456", p2: "123456", result: true},
		{p1: "123456", p2: "654321", result: false},
	} {
		hash, _ := HashPasswd(tt.p1)
		t.Log(hash)
		assert.Equal(t, tt.result, CheckPasswdHash(tt.p2, hash))
	}
}
