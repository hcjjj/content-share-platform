// Package web -----------------------------
// @file      : user_test.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-09 19:45
// -------------------------------------------
package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncrypt(t *testing.T) {
	password := "hello#worldhcj"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}
