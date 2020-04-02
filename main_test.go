package kansdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	AccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	SecretKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	kan, err := NewClient(AccessKey, SecretKey)
	if err != nil {
		panic(err)
	}

	err = kan.Email("topic", "msg")
	assert.Equal(t, "User not Exist", err.Error())
}
