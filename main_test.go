package caresdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	AccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	SecretKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	care, err := NewClient(AccessKey, SecretKey)
	if err != nil {
		panic(err)
	}

	err = care.Email("topic", "msg")
	//goland:noinspection GoNilness
	assert.Equal(t, "User not Exist", err.Error())
}

func TestLogPub(t *testing.T) {
	AccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	SecretKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	_, err := NewLogClient(AccessKey, SecretKey, "Log Topic", false)
	//goland:noinspection GoNilness
	assert.Equal(t, "User not Exist", err.Error())
}
