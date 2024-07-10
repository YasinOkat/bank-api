package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	acc, err := CreateAccount("a", "b", "yasin")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}
