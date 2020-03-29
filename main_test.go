package kan_sdk

import (
	"testing"
)

func TestSendEmail(t *testing.T) {
	AccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	SecretKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	
	kan, err := NewClient(AccessKey, SecretKey)
	if err != nil {
		panic(err)
	}

	kan.Email("topic", "msg")
}