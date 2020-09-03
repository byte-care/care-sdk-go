# care-sdk-go

# Install Care SDK
```bash
go get github.com/byte-care/care-sdk-go
```

# Quick Start
```go
package main
import "github.com/byte-care/care-sdk-go"

# Get AccessKey and SecretKey from https://www.bytecare.xyz/
AccessKey := "XXXXXX"
SecretKey := "XXXXXX"

care, err := NewClient(AccessKey, SecretKey)
if err != nil {
	panic(err)
}

care.Email("topic", "msg")
```