package kan_sdk

import (
	"bytes"
	"errors"
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

func (client *Client) consPostData(specificParameter map[string]string) (data url.Values) {
	commonParameter := sign.CommonParameter{
		client.credential.AccessKey,
		"sdfsdf",
		"4242",
	}

	signature := client.credential.Sign(commonParameter, specificParameter)

	data = map[string][]string{
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

	return
}

func (client *Client) post(specificParameter map[string]string) (err error) {
	data := client.consPostData(specificParameter)
	body := strings.NewReader(data.Encode())

	resp, err := http.Post("https://api.kan-fun.com/send-email", "application/x-www-form-urlencoded", body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return errors.New(buf.String())
	}

	return
}

func (client *Client) Email(topic string, msg string) (err error) {
	specificParameter := map[string]string{
		"topic": topic,
		"msg":   msg,
	}

	return client.post(specificParameter)
}
