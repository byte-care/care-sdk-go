# kan-sdk-go

# Install Kan SDK
```bash
go get github.com/kan-fun/kan-go-python
```

# Quick Start
```go
import "github.com/kan-fun/kan-go-python"

# Get AccessKey and SecretKey from http://www.mlflow.org.cn/
AccessKey := "XXXXXX"
SecretKey := "XXXXXX"

kan, err := NewClient(AccessKey, SecretKey)
if err != nil {
	panic(err)
}

kan.Email("topic", "msg")
```