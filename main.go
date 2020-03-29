package kan_sdk

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/kan-fun/kan-core"
)

type Client struct {
	credential *sign.Credential
}

func NewClient(accessKey, secretKey string) (client *Client, err error) {
	credential, err := sign.NewCredential(accessKey, secretKey)
	if err != nil {
		return
	}

	client = &Client{
		credential,
	}

	return
}

func post(url string, data url.Values) (resp *http.Response, err error) {
	body := strings.NewReader(data.Encode())

	resp, err = http.Post(url, "application/x-www-form-urlencoded", body)

	return
}

func (client *Client) Email(topic string, msg string) (err error) {
	commonParameter := sign.CommonParameter{
		client.credential.AccessKey,
		"sdfsdf",
		"4242",
	}

	specificParameter := map[string]string{
		"topic": topic,
		"msg":   msg,
	}

	signature := client.credential.Sign(commonParameter, specificParameter)

	data := map[string][]string{
		"access_key":      {commonParameter.AccessKey},
		"signature_nonce": {commonParameter.SignatureNonce},
		"timestamp":       {commonParameter.Timestamp},
		"signature":       {signature},
	}
	for k, v := range specificParameter {
		s := make([]string, 1)
		s[0] = v

		data[k] = s
	}

	_, err = post("http://58b5dd3da8514f30a8dfbf42bb0a740c-cn-beijing.alicloudapi.com/send-email", data)

	return
}
